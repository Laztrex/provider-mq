import asyncio
from pydantic import BaseModel

from framework.interfaces.base.api.base_model_pb2_grpc import BaseModelInterfaceServicer
from framework.interfaces.base.api.base_model_pb2 import ModelInterfaceResponse
from fastapi import BackgroundTasks

_REQUEST_ID_HEADER = 'request-id'


class Message(BaseModel):
    data: bytes
    headers: dict

    async def body(self):
        return self.data


class BaseModelGrpcInterface(BaseModelInterfaceServicer):

    def __init__(self, handler):
        """
        Init an instance of the service with the specified model
        as a request handler.

        :param handler: model instance.
        """

        self._model_predict = handler.predict
        self._model = handler()
        self._model._request_state = 2

    def Predict(self, request, context):
        """ Make prediction based on initial data in request. """
        request_id = dict(context.invocation_metadata()).get(_REQUEST_ID_HEADER.lower(), '')
        msg = Message(data=request.data, headers={'RqUID': request_id})
        answer = ModelInterfaceResponse(
            # data=asyncio.get_event_loop().run_until_complete(self._model.predict(request.data, request_id=request_id))
            data=asyncio.new_event_loop().run_until_complete(self._model_predict(
                self=self._model, item=msg, background_tasks=BackgroundTasks()),
            ))

        return answer
