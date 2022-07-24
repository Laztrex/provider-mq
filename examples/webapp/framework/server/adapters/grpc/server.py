import logging
import threading
from concurrent.futures import ThreadPoolExecutor
from secrets import randbelow
from typing import Callable, Any

from grpc import server, ServicerContext
from grpc_interceptor import ServerInterceptor, ExceptionToStatusInterceptor

from configs.app_config import CONFIGS

PORT = int(CONFIGS['server']['grpc']['port']) or (10000 + randbelow(5000))
_GRACE_PERIOD = int(CONFIGS['server']['grpc']['grace_period'])
_MAX_WORKERS = int(CONFIGS['server']['grpc'].get('max_server_pool'))
_REQUEST_ID_HEADER_NAME = 'RqUID'

_LOGGER = logging.getLogger(__name__)


class TracingServerInterceptor(ServerInterceptor):

    def intercept(self, method: Callable, request: Any, context: ServicerContext, method_name: str) -> Any:
        metadata = dict(context.invocation_metadata())
        request_id = metadata.get(_REQUEST_ID_HEADER_NAME.lower(), '')
        context.send_initial_metadata(((_REQUEST_ID_HEADER_NAME.lower(), request_id),))

        return method(request, context)


_server = server(
    ThreadPoolExecutor(max_workers=_MAX_WORKERS),
    interceptors=(TracingServerInterceptor(), ExceptionToStatusInterceptor())
)


def get_server():
    return _server


def shutdown():
    """ Stop gRPC server. """
    _LOGGER.info('Shutting down gRPC server..')
    _server.stop(_GRACE_PERIOD)


def run():
    """ Start gRPC server. """
    def start_server():
        _server.add_insecure_port(f'[::]:{PORT}')
        _server.start()

        try:
            _server.wait_for_termination()
        finally:
            _LOGGER.info('gRPC server stopped')

    threading.Thread(target=start_server, name='gRPCServerThread').start()
    _LOGGER.info(f'gRPC server listening on {PORT}..')
