import threading
from logging import Filter, LogRecord

from typing import Set

_local_context = threading.local()


def put(context: dict):
    """
    Put values into thread local context.
    :param context: key-value pairs to put into context.
    """
    for key, value in context.items():
        setattr(_local_context, key, value)


def get(key: str) -> object:
    """
    Get value from context.
    :param key: key.
    :return: Value from thread local context.
    """
    return getattr(_local_context, key, '')


def remove(keys: Set[str]):
    """
    Remove keys from thread local context.
    :param keys: keys to remove.
    """
    for key in keys:
        delattr(_local_context, key)


def contains(key: str) -> bool:
    return hasattr(_local_context, key)


class ThreadContextFilter(Filter):

    def __init__(self, filter_keys: Set[str] = None):
        super().__init__()
        self.__keys = filter_keys

    def filter(self, record: LogRecord):
        if self.__keys:
            for key in self.__keys:
                setattr(record, key, getattr(_local_context, key, ''))

        return True
