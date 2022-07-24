import json
import logging
import pika
import ssl
from abc import abstractmethod, ABCMeta
from typing import Optional, NoReturn

from aio_pika import connect_robust, ExchangeType, IncomingMessage
from aio_pika.abc import AbstractRobustChannel, AbstractRobustConnection, AbstractRobustQueue
from pika import BlockingConnection
from pika.adapters.blocking_connection import BlockingChannel

from configs.app_config import CONFIGS
from framework.schemas.body import MQMessage


class MQModel(metaclass=ABCMeta):
    def __init__(self):

        logging.getLogger("pika").setLevel(logging.WARNING)

        self._configs = CONFIGS.get("amqp")

        self._credentials = pika.PlainCredentials(
            self._configs.get("user"), self._configs.get("password")
        )

        self.connection: Optional[BlockingConnection] = None
        self.channel: Optional[BlockingChannel] = None
        self.response = None
        self.ssl_mode: bool = False

        if self._configs.get("tls"):
            context = ssl.create_default_context(cafile=self._configs["tls"].get("cafile"))
            context.load_cert_chain(self._configs["tls"].get("certfile"),
                                    self._configs["tls"].get("keyfile"))

            self.ssl_options = pika.SSLOptions(context, 'rabbitmq')
            self.ssl_mode = True

        self.connect(ssl_mode=self.ssl_mode)

        logging.info('Pika connection initialized')

    @abstractmethod
    async def predict(self, *args, **kwargs):
        """
        Основной метод для абстрактного класса, описывается в app.py
        """
        return NotImplementedError("Subclasses should implement this")

    @property
    def config(self) -> dict:
        return self._configs

    def connect(self, ssl_mode: bool = False) -> NoReturn:
        """
        Установление соединения для "ответного" канала. Используется в методе _publish_data
        :param ssl_mode: флаг ssl-соединения
        """

        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host=self.config.get("host"), port=self.config.get("port"),
                                      ssl_options=self.ssl_options if ssl_mode else None,
                                      credentials=self._credentials,
                                      heartbeat=60)
        )
        self.channel = self.connection.channel()

    async def declare_topic(self, channel: AbstractRobustChannel) -> AbstractRobustQueue:
        """
        Объявление топика, очереди и привязка этой очереди к топику
        :param channel: канал AbstractRobustConnection
        :return: Объявленная очередь
        """

        # topic_exchange = await channel.declare_exchange(self.config.get("TOPIC"), ExchangeType.TOPIC, durable=True,
        #                                                 passive=True)
        # queue = await channel.declare_queue(self.config.get("consumeQueue"), passive=True)

        # topic_exchange = await channel.get_exchange(self.config.get("TOPIC"))
        queue = await channel.get_queue(self.config.get("consumeQueue"))
        # await queue.bind(topic_exchange, routing_key=os.getenv("ROUTING_KEY"))

        return queue

    async def declare_queue(self, channel: AbstractRobustChannel):
        """
        Объявление очереди
        :param channel: канал AbstractRobustConnection
        :return: Объявленная очередь
        """
        queue = await channel.declare_queue(self.config.get("consumeQueue"))

        return queue

    async def get_queue(self, channel: AbstractRobustChannel) -> AbstractRobustQueue:
        """
        Получение существующей очереди
        :param channel: канал AbstractRobustConnection
        :return: Выделенная очередь
        """

        queue = channel.get_queue(self.config.get("consumeQueue"))

        return queue

    async def consume(self, loop) -> AbstractRobustConnection:
        """
        Ожидание сообщений в запущенном цикле событий
        :param loop: цикл событий, объявлен в app_server.py
        :return: Установленное соединение
        """

        connection = await connect_robust(
            host=self.config.get("host"),
            port=self.config.get("port"),  # 5672 - tcp, 5671 - tls
            loop=loop,
            login=self._credentials.username,
            password=self._credentials.password,
            ssl=True,
            ssl_options={
                "cafile": self.config["tls"].get("cafile"),
                "certfile": self.config["tls"].get("certfile"),
                "keyfile": self.config["tls"].get("keyfile")
            } if self.ssl_mode else None
        )
        channel = await connection.channel()

        # await channel.set_qos(prefetch_count=1)

        if self.config.get("TOPIC"):
            queue = await self.declare_topic(channel)
        else:
            queue = await self.declare_queue(channel)
            # queue = await self.get_queue(self.config.get("consumeQueue"))

        await queue.consume(self.predict, no_ack=False)
        logging.info('Established pika async listener')

        return connection

    async def get_data(self, message, *args, **kwargs) -> IncomingMessage:
        """Processing incoming message from RabbitMQ"""

        return message

    async def give_data(self, message: MQMessage, *args, **kwargs) -> NoReturn:
        """
        Публикация сообщений в RabbitMQ
        :param message: Сообщение, содержащее ответ ML модели
        """

        self.connect(ssl_mode=True)

        # try:

        await self._publish_data(message.result, message.headers,
                                 correlation_id=message.correlation_id,
                                 reply_to=message.reply_to)
        # except Exception as e:
        #     logging.info("Reconnection")
        #     self.connect(ssl_mode=True)
        #     await self._publish_data(message.result, message.headers,
        #                              correlation_id=message.correlation_id,
        #                              reply_to=message.reply_to)
        self.connection.close()

    async def _publish_data(
            self,
            message: dict, headers: dict,
            correlation_id: str = None, reply_to: str = None
    ) -> NoReturn:
        """Method to publish message to RabbitMQ"""
        logging.info(f"publish headers: {headers}")
        self.channel.basic_publish(
            exchange='',
            routing_key=reply_to,
            properties=pika.BasicProperties(
                headers=headers,
                correlation_id=correlation_id
            ),
            body=json.dumps(message, ensure_ascii=False).encode('utf-8'),
        )
