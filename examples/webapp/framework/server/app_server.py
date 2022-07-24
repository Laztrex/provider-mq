import importlib
import logging

# from multiprocessing import Manager
# from multiprocessing.managers import BaseProxy
from typing import Dict, List, NoReturn, Callable
from uuid import UUID

import uvicorn
from fastapi import FastAPI
from fastapi_utils.cbv import cbv
from fastapi_utils.inferring_router import InferringRouter

from configs.app_config import CONFIGS
from framework.server.contexts.process_context import InitProcessPool
from framework.server.handlers.health_check import health_checker
from framework.server.tech.release_resources import clear_workers
from framework.schemas.responses import PrettyJSONResponse
from framework.schemas.validation import Job
from model.initial import initializer_for_serv

app = FastAPI(title="ML-SERVING TEMPLATE",
              default_response_class=PrettyJSONResponse)

router = InferringRouter()

JOBS: Dict[UUID, Job] = {}

grpc_enabled: bool = CONFIGS['server'].get('grpc') is not None


##############################
#     CONFIGURE UVICORN      #
##############################


def start_server(Model):
    """
    Initial and start server with configs from *configs* or app_config.yml file (/configs/resources).
    HTTP or gRPC depending of the *configs*

    : configs : dict configuration server
                "server": {"uvicorn": {...},
                         "hypercorn": {...}}
                OR
                "grpc": {...}
    """

    initializer_for_serv()

    config_server = CONFIGS['server']
    config_workers = CONFIGS['framework'].get('workers')

    if not config_workers:
        raise Exception("Please, check your config-file *app_configs.yml* [field 'framework.workers' expected]")

    if grpc_enabled:
        start_grpc_server(Model, config_workers.get('grpc', {}))

    uvicorn.run(
        'framework.server.app_server:app',
        **config_server['uvicorn'],
        timeout_keep_alive=config_server.get('grace_period', 5)
    )

    if grpc_enabled:
        stop_grpc_server()


##############################
#      CONFIGURE gRPC        #
##############################


def start_grpc_server(model, config_workers: dict = None):
    """
    Инициализация grpc-сервера
    :param model: Класс используемой модели (из app.py)
    :param config_workers: параметры воркеров
    """
    from framework.interfaces.base.api import base_model_pb2_grpc as grpc_standard_model
    from framework.server.adapters.grpc import server as grpc_server
    from framework.interfaces.base.grpc.base_model_grpc import BaseModelGrpcInterface

    logging.info("Applying standard configuration to gRPC server..")

    initial_workers(**config_workers)

    standard_model_service = BaseModelGrpcInterface(model)
    server = grpc_server.get_server()
    grpc_standard_model.add_BaseModelInterfaceServicer_to_server(standard_model_service, server)
    grpc_server.run()


def stop_grpc_server():
    """Отключение grpc-сервера"""
    from framework.server.adapters.grpc import server as grpc_server

    logging.info("gRPC server shutdown..")

    grpc_server.shutdown()


##############################
#      INITIAL WORKERS       #
##############################


def initial_workers(
        max_workers: int = 1,
        timeout: float = None,
        initializer: dict = None,
        initargs: dict = None,
        *args, **kwargs
):
    """
    Инициализация воркеров.
    Передает параметры воркеров в класс InitProcessPool, который зависим от настроек в app_config.yml
    :param max_workers: макс. кол-во дочерних процессов (опционально, при None - по умолчанию инициализирующего класса
    [см. ProcessPoolExecutor])
    :param timeout: время работы
    :param initializer: функция инициализации воркеров
    :param initargs: аргументы инициилизирующей функции воркеров
    """

    # Initializer.DF_SHARED_VALUE = Manager().Namespace()

    if initializer:
        initializer: Callable = getattr(importlib.import_module(initializer.get('path')),
                                        initializer.get('name'))

        # Initializer.DF_SHARED_VALUE.df = pd.DataFrame(np.random.rand(100000, 20))
        # initargs = (Initializer.DF_SHARED_VALUE, )

        initargs = (initargs,)

    # TODO: add mapping executor's mode
    app.state.executor = InitProcessPool(
        max_workers=max_workers,
        timeout=timeout,
        initializer=initializer,
        initargs=initargs
    )


##############################
#     CONFIGURE ROUTES       #
##############################


def configure_routes() -> NoReturn:
    from app import Model

    router.add_api_route("/health", health_checker)
    router.add_api_route("/predict", getattr(Model, "predict"), methods=["POST"])

    check_route_background = getattr(Model, 'background', None)
    if check_route_background:
        router.add_api_route("/background", check_route_background, methods=["POST"])
        router.add_api_route("/status_job/{uid}", getattr(Model, 'status_job', None))

    # check_route_database = getattr(Model, 'in_database', None)
    # if check_route_database:
    #     router.add_api_route("/db", check_route_database, methods=["POST"])

    check_route_info = getattr(Model, 'info', None)
    if check_route_info:
        router.add_api_route("/info", check_route_info)

    check_route_cron = getattr(Model, 'cron', None)
    if check_route_info:
        router.add_api_route("/cron", check_route_cron)

    check_route_monitoring = getattr(Model, 'monitoring', None)
    if check_route_monitoring:
        router.add_api_route("/prometheus", check_route_monitoring)

    cbv(router)(Model)

    app.include_router(router)


##############################
#       EVENT STATE          #
##############################


@app.on_event("startup")
async def on_startup():
    """
    Starting FastAPI serv with configs from *app_configs.yml* for workers performance WORKERS.
    Decorated *on_event* is executed once.
    """

    initializer_for_serv()
    configure_routes()
    initial_workers(**CONFIGS['framework']['workers'].get('http', {}))


@app.on_event("shutdown")
async def on_shutdown(wait: bool = True):
    """
    Ending FastAPI serv, deactivating WORKERS.

    : wait : flag for mode deactivate executors;
            True - soft
            False - hard (manual process "deletion") and restart FastAPI serv (on_startup)
    """

    if wait:
        app.state.executor.shutdown(wait)
    else:
        clear_workers()
        JOBS.clear()
        await on_startup()
