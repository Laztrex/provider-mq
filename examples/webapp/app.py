#!/usr/bin/env python3

import logging

from collections import namedtuple
from typing import Union, List, NoReturn, Optional
from uuid import UUID

from fastapi import BackgroundTasks, Header, Depends
from starlette.requests import Request

from framework.interfaces.base.model import BaseModel
from framework.server.app_server import (
    JOBS, start_server
)
from framework.server.app_worker import start_model
from framework.interfaces.base.decorator import request_handler
from framework.schemas.validation import Job

from model import main


class Model(BaseModel):
    """
    Python -v 3.6+
    Example Service Model
    """

    __slots__ = ('_DATA',)

    @request_handler
    async def predict(
            self, item: Request,
            request_id: str = Header(''), content_type: str = Header(...)
    ) -> Union[dict, bytes]:
        """
        Основная функция старта расчета модели.

        Example full local request:
        >> DATA=$(echo -ne '{"data": [4, 15]}' | base64) (или из файла: DATA=$(openssl base64 -A -in ./01_refresh))
        >> RESP=$(curl -X POST -H "Content-Type: application/json" -H "RqUID: 37" -d '{"data": "'$DATA'"}' http://127.0.0.1:8080/predict)

        Прочитать 'сырое' сообщение
        >> echo -ne $RESP

        Декодированное сообщение
        >> python -c "import sys, json, base64; print(base64.b64decode($RESP['data']).decode('utf-8'))"

        :param item: Тело запроса
        :param request_id: Идентификатор запроса
        :param content_type: Тип запроса
        :return: Тело ответа
        """

        logging.info("Processing request started")
        logging.info(f"Request ID: {self.HEADERS.get('request-id')}")

        Params = namedtuple('Params', ['function', 'data'])
        worker = Params(main.start_predict, self.DATA)

        ans = await start_model(worker.function, None, worker.data)

        logging.info("Processing request done")

        return {"result": ans}

    @request_handler
    async def background(
            self, item: Request, background_tasks: BackgroundTasks,
            request_id: str = Header(''), content_type: str = Header(...)
    ) -> dict:
        """

        :param item: Тело запроса
        :param background_tasks: Очередь фоновых задач
        :param request_id: Идентификатор запроса
        :param content_type: Тип запроса
        :return: Тело ответа
        """
        Params = namedtuple('Params', ['function', 'data'])
        worker = Params(main.start_predict, self.DATA.get('data'))

        new_task = Job()
        logging.info(f"... to background with {new_task.uid}")
        JOBS[new_task.uid] = new_task
        background_tasks.add_task(start_model, worker.function, new_task.uid, worker.data)

        return new_task.asdict()

    async def status_job(self, uid: UUID) -> dict:
        """
        Проверка статуса задачи
        :param uid: Идентификатор задачи
        """

        result_job: dict = JOBS[uid].asdict()
        del JOBS[uid]

        return result_job


if __name__ == "__main__":
    start_server(Model)
