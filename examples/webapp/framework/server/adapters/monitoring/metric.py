from prometheus_client import Counter, Gauge, Summary

REQUEST_TAG = ['cluster_name', 'project_name', 'application_name', 'model_type', 'model_name', 'model_id']

RESPONSE_TAGS = REQUEST_TAG[:]
RESPONSE_TAGS.append('result')
RESPONSE_TAGS.append('exception')

MODEL_REQUEST_TOTAL = Counter('model_request_total', 'Всего количество запросов к модели со стороны АС', REQUEST_TAG)
MODEL_RESPONSE_TOTAL = Counter('model_response_total', 'Всего количество ответов модели на запрос полученный от АС', RESPONSE_TAGS)
MODEL_RESPONSE_SECONDS_SUMMARY = Summary('model_response_seconds', 'Время исполнения модели, в секундах', RESPONSE_TAGS)
MODEL_RESPONSE_SECONDS_GAUGE = Gauge('model_response_seconds_max', 'Время исполненич модели, в секундах', RESPONSE_TAGS)
