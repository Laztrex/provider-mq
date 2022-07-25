from dataclasses import dataclass
from typing import Optional, Union, Dict, List

from aio_pika.message import HeaderProxy
from pydantic import BaseModel


@dataclass
class FibBody:
    data: bytes

    async def body(self) -> str:
        return self.data.decode('utf-8')

    @property
    def headers(self, *args, **kwargs):
        new_headers = {'Content-Type': 'application/json'}
        new_headers.update(**kwargs)
        return new_headers


class RESTMessage(BaseModel):
    data: bytes
    _headers: dict = {'Content-Type': 'application/json'}

    async def body(self):
        return self.data

    @property
    def headers(self, *args, **kwargs):
        new_headers = self._headers
        new_headers.update(**kwargs)
        return new_headers


@dataclass
class MQMessage:
    data: Union[bytes, dict]
    RqUID: Optional[str] = None
    headers: Optional[dict] = None
    result: Optional[dict] = None

    _correlation_id: Optional[str] = None
    _reply_to: Optional[str] = None

    async def body(self):
        self.data = literal_eval(self.data.decode('utf-8'))
        return self

    @property
    def correlation_id(self):
        return self._correlation_id

    @property
    def reply_to(self):
        return self._reply_to

    @correlation_id.setter
    def correlation_id(self, value):
        self._correlation_id = value

    @reply_to.setter
    def reply_to(self, value):
        self._reply_to = value


@dataclass
class MQHeaders:
    _request_headers: Optional[HeaderProxy] = None

    def __getitem__(self, item):
        return self._request_headers.get(item)

    def get(self, item):
        return self._request_headers.get(item)