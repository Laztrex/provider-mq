import threading
from datetime import datetime

start_time = datetime.now()
shutdown_event = threading.Event()
