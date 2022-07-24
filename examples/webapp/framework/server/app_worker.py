import asyncio
import logging
from concurrent.futures.process import BrokenProcessPool

from typing import Coroutine, NoReturn, List
from uuid import UUID

from framework.server.app_server import app, JOBS, on_shutdown, on_startup


async def start_multiple_model(workers: List):
    workers_task = [
        start_model(worker.num_worker, worker.function, None, worker.event_id, worker.text, worker.role)
        for worker in workers
    ]
    results = await asyncio.gather(*workers_task)
    return results


async def start_model(function_worker, uid: UUID = None, *args):
    """
    Функция приёма задач и отправка в цикл событий
    :param function_worker:
    :param uid: идентификатор задачи
    """

    ans = await run_process(
        function_worker,
        *args
    )

    if uid:
        if not JOBS[uid].no_return:
            JOBS[uid].result = ans
        JOBS[uid].status = "complete"
    else:
        return ans


async def run_process(task, *args) -> Coroutine:
    """
    Инициализация цикла событий
    :param task: Функция точки входа в модель
    :param args: Параметры модели
    :return:
    """

    try:

        return await app.state.executor.run(task, *args)  # position-dependent arguments

    except (BrokenProcessPool, MemoryError) as exc:
        logging.info(f"{exc}\n Deficiency resources to processes the model.\n"
                     f"Try decreasing 'max_workers' or optimizing computations")

        await on_shutdown()
        await on_startup()


async def del_job(uid: UUID) -> NoReturn:
    """
    Удаление задачи
    :param uid: Идентификатор задачи
    """

    del JOBS[uid]
    logging.info(f'tech_status: Delete job with uid {uid}')


async def clear_all_jobs() -> NoReturn:
    """
    Удаление всех задач
    """
    JOBS.clear()
    logging.info(f'tech_status: All jobs is clear')


async def check_spec_fields(data: dict):
    """
    Проверка наличия в запросе специальных полей для работы воркерами / фоновыми задачами
    """

    if data.pop('clear_workers', None):
        await on_shutdown(False)
        logging.info('tech_status: Workers have been reloaded')

    if data.pop('clear_all_jobs', None):
        await clear_all_jobs()

    if data.get('del_job'):
        uid_job = data.pop('del_job')
        await del_job(
            UUID(uid_job)
        )
