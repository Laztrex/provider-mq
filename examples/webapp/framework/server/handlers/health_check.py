from framework.server.contexts.singleton import Singleton


def health_checker() -> dict:
    """
    Проверка работоспособности
    Example local request:
    curl http://127.0.0.1:8080/health
    """

    return {"health_status": "running"}


class HealthIndicator(metaclass=Singleton):
    """future"""
    def get_status(self, *args, **kwargs):
        pass
