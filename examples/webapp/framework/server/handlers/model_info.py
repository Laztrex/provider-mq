import logging

from configs.app_config import CONFIGS
from framework.server.contexts.singleton import Singleton
from framework.interfaces.base.decorator import technical_request_handler

_LOGGER = logging.getLogger(__name__)


class Info(metaclass=Singleton):
    """
    Singleton, representing model info.
    From app_config.yml file
    """

    def __init__(self):
        self._info = CONFIGS.get('model')

    @technical_request_handler
    def get_info(self, **context) -> dict:
        if context.get('shutdown_event').is_set():
            _LOGGER.warning("Shutdown process is started. Trying to respond")
        return self._info
