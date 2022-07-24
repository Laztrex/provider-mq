import logging
import os


def clear_workers():
    """
    Test func for del process
    """

    import signal
    import psutil

    def kill_child_processes(parent_pid, sig=signal.SIGTERM):
        try:
            logging.info(f'clear process')

            parent = psutil.Process(parent_pid)
        except psutil.NoSuchProcess:
            logging.info(f'No Such Process for {parent_pid}')
            return
        children = parent.children(recursive=True)
        for process in children:
            process.send_signal(sig)

    kill_child_processes(os.getpid())
