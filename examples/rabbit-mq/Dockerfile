FROM rabbitmq:3-management

# Define environment variables.
ENV RABBITMQ_ADMIN_USER superuser
ENV RABBITMQ_ADMIN_PASSWORD superpassword

ENV RABBITMQ_USER MLUser
ENV RABBITMQ_PASSWORD MLPass

ENV RABBITMQ_PID_FILE /var/lib/rabbitmq/mnesia/rabbitmq

ADD init.sh /init.sh
RUN chmod +x /init.sh

EXPOSE 15672

# Define default command
CMD ["/init.sh"]