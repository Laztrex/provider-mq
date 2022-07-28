FROM python:3

ENV CN=rabbitmq

RUN git clone https://github.com/michaelklishin/tls-gen tls-gen
RUN cd /tls-gen/basic && make PASSWORD= CN=${CN} && make verify

#RUN cat /tls-gen/basic/result/ca_certificate.pem
#RUN cat /tls-gen/basic/result/client_certificate.pem
#RUN cat /tls-gen/basic/result/client_key.pem
#RUN cat /tls-gen/basic/result/server_certificate.pem
#RUN cat /tls-gen/basic/result/server_key.pem