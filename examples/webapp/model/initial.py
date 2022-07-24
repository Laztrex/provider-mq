import logging.config
import sys
# from multiprocessing.managers import Namespace

from typing import NoReturn


class Initializer(object):
    """
    Инициализция параметров
    """

    # DF_SHARED_VALUE = None

    @staticmethod
    def init_logger_serv(disable_exist_log: bool = True) -> NoReturn:
        """
        Инициализация конфиграции логирования
        """
        from configs import logger_config_loader

        logger_config_loader.load(disable_exist_log)

    @staticmethod
    def init_logger_model(name) -> NoReturn:
        root = logging.getLogger(name)
        root.setLevel(logging.INFO)

        handler = logging.StreamHandler(sys.stdout)
        handler.setLevel(logging.INFO)
        formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
        handler.setFormatter(formatter)
        root.addHandler(handler)

        logging.info("Logging is set!")

    @classmethod
    def set_meta(cls, name: str = 'meta_data', value=None):
        setattr(cls, name, value)

    @staticmethod
    def get_log():
        return Initializer.init_logger_serv()


def initializer_for_serv(*args, **kwargs):
    Initializer.init_logger_serv(False)


def initializer_for_model(value, *args, **kwargs):
    """Initialize function for every worker [from framework.workers.http.initializer in app_config.yml]"""
    # Initializer.DF_SHARED_VALUE = value
    return Initializer.init_logger_serv(False)
