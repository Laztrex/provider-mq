from os import getcwd
from os.path import abspath, dirname, join, exists, pardir

import yaml

CONFIG_PATH = join(getcwd(), 'configs', 'resources')


def get_db_config(config_name: str = 'db_config.yml') -> list:
    """
    Чтение настроек для provider-db. В т.ч. ожидаемый порядок столбцов, типы
    :param config_name: Имя конфигурационного файла
    :return: Список конфигураций
    """

    config = join(CONFIG_PATH, config_name)
    if not exists(config):
        return []

    with open(config, 'r') as cfg:
        d = yaml.safe_load(cfg.read())

        list_db_config = []
        for k, v in d['tables'].items():
            if v.get('params'):
                v['params'].update(_check_enc_exec_data(v['params']))
            list_db_config.append({'query_id': k, **v})
        return list_db_config


def _check_enc_exec_data(to_exec_data: dict) -> dict:
    """
    Enc params for db
    :param to_exec_data: input params
    :return: encoded params
    """

    return dict((k, str(v).encode('utf-8')) for (k, v) in to_exec_data.items())
