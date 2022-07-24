import requests

from framework.server.contexts import thread_context

http = requests.Session()

http.headers.update({'RqUID': str(thread_context.get('request_id'))})
http.hooks['response'] = [lambda response, *args, **kwargs: response.raise_for_status()]
