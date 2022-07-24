import logging

from time import time

from configs.app_config import CONFIGS
from framework.server.adapters.monitoring.logger_metric import get_model_metric
from framework.server.adapters.monitoring.metric import MODEL_REQUEST_TOTAL, MODEL_RESPONSE_TOTAL, \
    MODEL_RESPONSE_SECONDS_SUMMARY, MODEL_RESPONSE_SECONDS_GAUGE

_LOGGER = logging.getLogger('metric')
_CLUSTER_NAME = CONFIGS['cluster']['name']
_PROJECT_NAME = CONFIGS['project']['name']
_APPLICATION_NAME = CONFIGS['application']['name']
_MODEL_TYPE = CONFIGS['model']['type']
_MODEL_NAME = CONFIGS['model']['name']
_MODEL_ID = CONFIGS['model']['id']


def start_metric() -> 'time':
    start_time = time()
    set_count_request_metric()
    return start_time


def finish_metric(start_time: 'time', status: str, error: str) -> None:
    spent = time() - start_time
    set_count_response_metric(status, error)
    set_timed_metric(spent, status, error)


def set_count_request_metric() -> None:
    MODEL_REQUEST_TOTAL.labels(
        cluster_name=_CLUSTER_NAME,
        project_name=_PROJECT_NAME,
        application_name=_APPLICATION_NAME,
        model_type=_MODEL_TYPE,
        model_name=_MODEL_NAME,
        model_id=_MODEL_ID
    ).inc()

    metric = get_model_metric('model_request_total', str(1), '', '')
    _LOGGER.info(metric)


def set_count_response_metric(result: str, error: str) -> None:
    MODEL_RESPONSE_TOTAL.labels(
        cluster_name=_CLUSTER_NAME,
        project_name=_PROJECT_NAME,
        application_name=_APPLICATION_NAME,
        model_type=_MODEL_TYPE,
        model_name=_MODEL_NAME,
        model_id=_MODEL_ID,
        result=result,
        exception=error
    ).inc()

    metric = get_model_metric('model_response_total', str(1), result, error)
    _LOGGER.info(metric)


def set_timed_metric(spent: 'time', result: str, error: str) -> None:
    MODEL_RESPONSE_SECONDS_SUMMARY.labels(
        cluster_name=_CLUSTER_NAME,
        project_name=_PROJECT_NAME,
        application_name=_APPLICATION_NAME,
        model_type=_MODEL_TYPE,
        model_name=_MODEL_NAME,
        model_id=_MODEL_ID,
        result=result,
        exception=error
    ).observe(spent)

    MODEL_RESPONSE_SECONDS_GAUGE.labels(
        cluster_name=_CLUSTER_NAME,
        project_name=_PROJECT_NAME,
        application_name=_APPLICATION_NAME,
        model_type=_MODEL_TYPE,
        model_name=_MODEL_NAME,
        model_id=_MODEL_ID,
        result=result,
        exception=error
    ).set(spent)

    metric = get_model_metric('model_response_seconds', str(spent), result, error)
    _LOGGER.info(metric)
