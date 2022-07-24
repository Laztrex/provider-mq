__MODEL_TAG_MAP = "result=\"{}\", exception=\"{}\""


def get_model_metric(
        metric_name: str,
        value: str,
        result: str,
        exception: str
) -> str:
    tag = __MODEL_TAG_MAP.format(result, exception)
    metric = "METRIC " + metric_name + " " + value + " {" + tag + "}"
    return metric
