from functools import wraps


def mq_handler(fn):
    @wraps(fn)
    async def wrapper(self, message, *args, **kwargs):
        from framework.schemas.body import MQMessage, MQHeaders

        incoming_msg = MQMessage(data=message.body)
        incoming_headers = MQHeaders(_request_headers=message.headers)

        RqUID = incoming_headers.get('RqUID')
        content_type = incoming_headers.get('Content-Type')
        incoming_msg.correlation_id = message.correlation_id
        incoming_msg.reply_to = message.reply_to
        incoming_msg.RqUID = RqUID

        incoming_msg.headers = {"RqUID": RqUID}

        result = await fn(
            self,
            item=incoming_msg,
            RqUID=RqUID,
            content_type=content_type,
            headers=incoming_headers,
            *args, **kwargs
        )

        await message.ack()

        return result

    return wrapper
