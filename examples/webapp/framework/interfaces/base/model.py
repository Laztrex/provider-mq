import json
import logging
from abc import abstractmethod, ABCMeta
from ast import literal_eval
from base64 import b64encode
from typing import Union, NoReturn, Optional
from uuid import UUID

from google.protobuf.json_format import Parse
from starlette.responses import Response

from configs.app_config import CONFIGS
from framework.interfaces.base.api.base_model_pb2 import ModelInterfaceRequest
from framework.schemas.responses import PrettyJSONResponse


class BaseModel(metaclass=ABCMeta):
    """
    Example Base Service Model with two main routers: /health, /predict

    The class attribute type affects the REST request.
    If the expected data type is different from dict, override the type in the inherited "Model" clas (app.py).
    For example: self.DATA: dict = None
    """

    __slots__ = ('DATA', 'HEADERS', '_request_mapping', '_response_mapping', '_request_state')

    _BACKGROUND = CONFIGS['framework'].get('background_tasks')

    def __init__(self):
        self.DATA: Union[dict, bytes] = {}
        self.HEADERS: dict = {}
        self._request_mapping = {
            1: self._base_request, 2: self._base_b64_request, 3: self._grpc_request
        }
        self._response_mapping = {
            1: self._base_response, 2: self._base_b64_response, 3: self._media_response
        }
        self._request_state: int = 1

    @abstractmethod
    async def predict(self, *args, **kwargs):
        return NotImplementedError("Subclasses should implement this")

    @property
    def __request(self):
        """Mapping parse-method request body"""
        return self._request_mapping[self._request_state]

    @property
    def __response(self):
        """Mapping parse-method request body"""
        return self._response_mapping[self._request_state]

    def _base_request(self, data: str, *args, **kwargs):
        """
        Mock deserialize while use dataclass body in /schemas.
        Compatibility with decorator
        """

        return literal_eval(data)

    def _base_b64_request(self, data: bytes, content_type: str) -> dict:
        """Deserialize request data [http]"""

        if content_type == "application/json":
            parsed_request = Parse(data, ModelInterfaceRequest())
            parsed_data = parsed_request.data.decode('utf-8')
        else:
            parsed_request = ModelInterfaceRequest()
            parsed_request.ParseFromString(data)
            parsed_data = parsed_request.data.decode('utf-8')

        return literal_eval(parsed_data)

    def _grpc_request(self, data: bytes) -> dict:
        return literal_eval(data.decode('utf-8'))

    def _base_response(self, data: Union, *args, **kwargs) -> Union[PrettyJSONResponse, Response]:
        """
        Method for response message layout [http]
        :param data: output body message
        """

        if not data:
            return Response(status_code=200)

        return PrettyJSONResponse(
            content=data,
            headers=self.HEADERS,
            media_type="application/json"
        )

    def _base_b64_response(self, data: Union, *args, **kwargs) -> Union[PrettyJSONResponse, Response]:
        """
        Method for response message layout [http]
        :param data: output body message
        """

        if not data:
            return Response(status_code=200)

        content = {"data": b64encode(bytes(json.dumps(data, ensure_ascii=False), 'utf-8'))}

        return PrettyJSONResponse(
            content=content,
            headers=self.HEADERS,
            media_type="application/json"
        )

    def _media_response(self, data, content_type: str, *args, **kwargs) -> Optional:
        """
        Method for response message layout
        :param data: output media body message
        """

        self._request_state = 1

        if all(
                [not isinstance(data, Response), not isinstance(data, dict)]
        ):
            return Response(content=data, headers=self.HEADERS, media_type=content_type)
        else:
            return data

    def _grpc_response(self, data: Union, *args, **kwargs):
        """
        Method for response message layout
        :param data: output body message
        """

        return json.dumps(data).encode('utf-8')

    async def get_data(self, item: Union[bytes, dict], headers: dict, request_id: str = '', content_type: str = None,
                       *args, **kwargs) -> NoReturn:
        """
        Parser request body.
        Include with protobuf protocol
        :param item: Input message from request
        :param headers: headers from request
        :param request_id: Request identification
        :param content_type: type request body
        """

        self.DATA = self.__request(item, content_type)

        self.HEADERS.update({"request-id": headers.get("request-id", request_id)})

    async def give_data(self, data, *args, **kwargs):
        """
        Check validate Response class
        :param data: Response body message
        """
        return self.__response(data, *args, **kwargs)
