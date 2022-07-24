import logging

from requests import Response

from configs.app_config import CONFIGS
from .http import http

from model.model_app.api.generated.controller_storage_pb2 import StorageCommand

from ..handlers import locator

controller_storage_name = CONFIGS['application']['controller_storage_name']

_LOGGER = logging.getLogger(__name__)


def use_storage(command: int, s3url: str) -> Response:
    if not s3url:
        raise ValueError('s3url must not be None or empty')

    if (command not in StorageCommand.values()):
        error_msg = 'Invalid command, only allowed [{}]'.format(','.join(str(value) for value in StorageCommand.values()))
        raise ValueError(error_msg)

    body = {
        'command': command,
        's3Url': s3url
    }

    return post_request(body)


def post_request(body: str) -> Response:
    """POST request to controller storage.

    :param body: json data to send in body of Request.
    :return: controller storage Response.
    :rtype: requests.Response
    """

    url = 'http://' + locator.find_socket(controller_storage_name) + '/v1/api/use-storage'
    headers = {'Content-Type': 'application/json'}

    return http.post(url, json=body, headers=headers, stream=True)
