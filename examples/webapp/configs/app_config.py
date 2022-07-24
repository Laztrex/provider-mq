import importlib.machinery
import logging
from os import getenv
from os.path import abspath, dirname, join, exists, pardir

import yaml
from jinja2 import Template

from framework.server.contexts import app_context as context


def _read_metadata(path_metadata: str) -> dict:
    """
    Загрузка файла метаданных (__metadata__) из директории модели /model
    Используется для интерфейса ManagementModel
    :param path_metadata: Путь к фалу метаданных
    :return: Словарь с метаданными
    """

    metadata: dict = {}

    try:
        if exists(path_metadata):
            loaded_metadata = importlib.machinery.SourceFileLoader('__metadata__', path_metadata).load_module()
            all_attrs = getattr(loaded_metadata, '__all__')
            if all_attrs is not None:
                for attr in all_attrs:
                    metadata[attr] = getattr(loaded_metadata, attr)
        return metadata

    except Exception:
        logging.info(f"Failed load metadata: {path_metadata}")


def read_config() -> dict:
    framework_dir = abspath(join(__file__, pardir, pardir))
    model_dir = join(framework_dir, 'model')
    config_dir = join(dirname(__file__), '')

    config_file = join(config_dir, 'resources', 'app_config.yml')
    framework_metadata_file = join(framework_dir, '__metadata__.py')
    model_metadata_file = join(model_dir, '__metadata__.py')

    if not exists(config_file):
        return {}
    with open(config_file, 'r') as cfg:
        result = yaml.safe_load(
            Template(cfg.read()).render(
                env=getenv,
                app_context_start_time=context.start_time.isoformat(),
                model_metadata=_read_metadata(model_metadata_file),
                framework_metadata=_read_metadata(framework_metadata_file)
            )
        )

        return result


CONFIGS: dict = read_config()
