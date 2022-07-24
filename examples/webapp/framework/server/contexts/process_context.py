import asyncio
import logging
import multiprocessing

from concurrent.futures import ProcessPoolExecutor
from multiprocessing import BoundedSemaphore
from threading import BoundedSemaphore
from typing import Callable, Iterable


class InitProcessPool:
    """
    This class is used to initialize the executor for the workers.
        * For <timeout> or <cancellable> parameters, an instance of the "CancellablePool" class is used.
                "State-method": _run_cancellable.
        * With only <max_workers> parameter and "concurrent.futures.ProcessPoolExecutor" is used.
                "State-method": _run_concurrent.
        * In the absence of any parameters -
                 "workers" are executed in the event loop in the current process (blocking it).
                 "State-method": _run_concurrent.
    You can transfer "your" executor.
    """

    def __init__(self, max_workers: int = None,
                 timeout: float = None,
                 executor=None,
                 cancellable: bool = False,
                 initializer: Callable = None,
                 initargs: Iterable = None,
                 msg: str = '',
                 *args, **kwargs):

        __MAX_WORKERS = max_workers if max_workers else 1

        if executor:
            self.executor = executor(*args, **kwargs)
            self.__state = getattr(self, '_run_concurrent')

        elif timeout or cancellable:
            self.TIMEOUT = timeout
            self.executor = CancellablePool(max_workers=__MAX_WORKERS, initializer=initializer, initargs=initargs)
            self.__state = getattr(self, '_run_cancellable')
            self.futures = asyncio.Queue()

        else:
            self.executor = ProcessPoolExecutor(max_workers=__MAX_WORKERS)  # add initializer for python 3.7+
            self.__state = getattr(self, '_run_concurrent')

        logging.info(f"Connected {__MAX_WORKERS} executor's: "
                     f"{self.executor.__class__.__name__}, state: {self.__state.__name__}: {msg}")
        # or MaxQueuePool(ProcessPoolExecutor(**init_params)) in ./context/process_context.py
        # or from pebble import ProcessPool

    async def _run_cancellable(self, loop, task, *args):
        """
        Arrange for func to be called in the "CancellablePool" executor.
        """

        self.futures.put_nowait(loop.create_task(self.executor.submit(task, *args)))
        object_new = self.futures.get_nowait()

        try:

            answer = await asyncio.wait_for(object_new, timeout=self.TIMEOUT)

        except asyncio.TimeoutError:
            answer = 'task timed out and cancelled'
            logging.info(answer)

        return answer

    async def _run_concurrent(self, loop, task, *args):
        """
        Arrange for func to be called in the "concurrent.futures.ProcessPoolExecutor" executor or standalone.
        """

        return await loop.run_in_executor(self.executor, task, *args)

    async def run(self, *args, **kwargs):
        """
        Runs "executor" in the event loop for execution according to "State-methods".
        """

        loop = asyncio.get_event_loop()
        return await self.__state(loop, *args, **kwargs)

    def shutdown(self, *args, **kwargs):
        return self.executor.shutdown(*args, **kwargs)


class CancellablePool:
    """
    This class creates a specified number of "workers" and ""terminate" them after execution.
    Used when app_configs.yml is set to "timeout" or "calcellable".
    Warning! May bring additional costs for imports. Unlike ProcessPool, they are produced immediately.
    """

    def __init__(self, max_workers: int = 3, initializer: Callable = None, initargs: Iterable = None):

        self.initializer = initializer
        self.initargs = initargs

        self._free = {self._new_pool() for _ in range(max_workers)}
        self._working = set()
        self._change = asyncio.Event()

    def _new_pool(self):
        """Create new *worker*"""
        return multiprocessing.Pool(1, initializer=self.initializer, initargs=self.initargs)

    async def submit(self, fn, *args):
        """
        Like multiprocessing.Pool.apply_async, but:
         * is an asyncio coroutine
         * terminates the process if cancelled
        """

        while not self._free:
            await self._change.wait()
            self._change.clear()
        pool = usable_pool = self._free.pop()
        self._working.add(pool)

        loop = asyncio.get_event_loop()
        fut = loop.create_future()

        def _on_done(obj):
            loop.call_soon_threadsafe(fut.set_result, obj)

        def _on_err(err):
            loop.call_soon_threadsafe(fut.set_exception, err)

        pool.apply_async(fn, args, callback=_on_done, error_callback=_on_err)

        try:
            return await fut
        except asyncio.CancelledError:
            pool.terminate()
            usable_pool = self._new_pool()
        finally:
            self._working.remove(pool)
            self._free.add(usable_pool)
            self._change.set()

    def shutdown(self, wait=None):
        for p in self._working | self._free:
            p.terminate()
        self._free.clear()


class MaxQueuePool:
    """
    This class wraps a concurrent.futures.Executor
    limiting the size of its task queue.
    If 'max_queue_size' tasks are submitted, the next call to submit will block
    until a previously submitted one is completed.
    """

    def __init__(self, executor, max_queue_size, max_workers=None):
        self.pool = executor(max_workers=max_workers)
        self.pool_queue = BoundedSemaphore(max_queue_size)

    def submit(self, function, *args, **kwargs):
        """Submits a new task to the pool, blocks if Pool queue is full"""

        self.pool_queue.acquire()

        future = self.pool.submit(function, *args, **kwargs)
        future.add_done_callback(self.pool_queue_callback)

        return future

    def pool_queue_callback(self, _):
        """Called once task is done, release one queue slot"""

        self.pool_queue.release()
