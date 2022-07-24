from logging import config, Filter, LogRecord
from os import getenv
from os.path import join, dirname

from jinja2 import Template
from yaml import safe_load

config_dir = join(dirname(__file__), '')
config_file = join(config_dir, 'resources', 'logger_config.yml')


class HealthCheckFilter(Filter):
    def filter(self, record: LogRecord) -> bool:
        return record.getMessage().find("/health") == -1


def load(disable_exist_log: bool = True):
    with open(config_file, 'r') as cfg_file:
        config.dictConfig(
            safe_load(
                Template(cfg_file.read()).render(
                    env=getenv,
                    disable_exist=disable_exist_log
                )
            )
        )
