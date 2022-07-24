# Provider-MQ

The project is part of a cloud-based microservice solution for interacting with an ML-application and related
infrastructure components.  


This application allows you to call for execution an ML model located in a dedicated server.   
*Consumer* is inside the application and when a message appears, a REST-request is made in the application with the ML model.

## Description

*Prover-MQ* currently implements the following:

* Publishing - REST request to ML model service
* Defining queue parameters in a configuration file
* Application arranges sync/async calls with the model
* Response from the model (sync)
    * If the pod with the model crashes before being called or during the calculation - the message will not be lost
      under certain conditions
* Declare of exchange, queues on the *provider-mq* side - we do not trust the model

The project will be updated.

## Usage

The project uses global environment variables.

| Name         |           Goal           |         Expected value example         |
|--------------|:------------------------:|:--------------------------------------:|
| RMQ_URL      | Host to connect RabbitMQ | "amqps://user:password@rabbitmq:5671/" |
| MQ_CACERT    |            CA            |   "/certs/provider_mq/cacert_mq.pem"   |
| MQ_CERT      |       Certificate        | "/certs/provider_mq/client_cert_mq.pem" |
| MQ_KEY       |     Key certificate      | "/certs/provider_mq/client_key_mq.pem"  |
| GIN_CERT     |    Server certificate    |  "/certs/provider_mq/server_cert.pem"   |
| GIN_CERT_KEY |        Server key        |   "/certs/provider_mq/server_key.pem"   |
| LOG_LEVEL    |      Logging level       |                "debug"                 |

You can define the default values in the **.env** file of the project root.

Configuration file is also provided to define protocol settings - [queue_config.yml](https://github.com/Laztrex/provider-mq/blob/master/configs/)

~~~yaml
- topic: ML.MQ
  queueName: fib
  bindingKey: "predict.*"
  replyTo: "response"
  dle: true
  argQueue: {
    "x-dead-letter-exchange": dead_letter_exchange,
  }
- topic: ML.MQ
  queueName: manage
  bindingKey: "metrics.*"
  replyTo: "response"

~~~

Currently, only Topic can be defined from the configuration file for Exchange (but Direct can also be defined). Other exchanges - will be supplemented.
Flexible settings for the queue - lifetime, autodelete, types, arguments and etc. currently not included in the configuration file, the parameters can be configured inside the code optionally.  

The directory [examples/webapp](https://github.com/Laztrex/provider-mq/blob/master/examples/webapp/) contains a simple web application for testing the project.

~~~
>> go version
go version go1.18.2 darwin/amd64
~~~

### Compose it

The project provides examples of services for testing provider-mq.

~~~
docker-compose build --no-cache
~~~

~~~
docker-compose up -d
~~~

~~~
docker-compose logs -f -t
~~~

Проверим:

~~~
>> examples % curl -k --key certs/client/client_key_mq.pem --cert certs/client/client_cert_mq.pem https://0.0.0.0:5050/v1/predict -d '{"data": "[15, 29]"}' -H "RqUID: 52-42" -H "Content-Type: application/json" -H "routing-key: predict.online"
~~~

## Objective

The goal of the project is to create binding services for the ML model that provide ease of interaction, flexible logic,
scalability of connecting various interfaces of the ML Engine architecture - *Cloud-based Sandbox for ML Serving*

Simplified sketch  
![Image alt](https://github.com/Laztrex/provider-mq/blob/master/docs/pics/first_sketch.png)

## Addition

It is implied in the microservice architecture of a cloud solution for calculating ML models:

* Prefix "**provider**-x-x" - auxiliary controllers to implement the interface with the infrastructure technical stack
* Prefix "x-**gateway**-x" - optional label to identify the converter, here indicates the presence of a REST-MQ adapter
* Predix "x-x-**mq**" - integrated interface

### Certificates

In this project, we are going to create a Golang web client to connect to RabbitMQ server with TLS. For this we
you will need to create self-signed SSL certificates and share them between the Golang application and the RabbitMQ server.

Directory [examples/certs](https://github.com/Laztrex/provider-mq/blob/master/examples/certs/) contains a Dockerfile with an example of generating a self-signed certificate.

~~~
cd examples/certs/
~~~
~~~
docker build -t certs .
~~~
~~~
docker run -i -t certs bash
>> cd /tls-gen/basic/result
~~~

Add certificate files to directory [examples/certs](https://github.com/Laztrex/provider-mq/blob/master/examples/certs/). Initially set certificate structure:

~~~
├── examples  
│   └── certs 
│       ├── rabbitmq
│       │    ├── cacert.pem  
│       │    ├── server_cert.pem  
│       │    └── server_key.pem  
│       ├── provider_mq
│       │    ├── cacert.pem  
│       │    ├── client_cert.pem  
│       │    └── client_key.pem 
│       └── webapp
│            ├── cacert.pem  
│            ├── client_cert.pem  
│            └── client_key.pem   
└── docker-compose.yaml  
~~~

It is also necessary to provide for the creation of certificates for the *gin*-server.

### Initial setup RabbitMQ

When initializing the RabbitMQ client, you can set initial settings.  
In particular, you can set available accounts and define a configuration file. See example in [examples/rabbit-mq](https://github.com/Laztrex/provider-mq/blob/master/examples/rabbit-mq/).  

Check the list of available users:

~~~
>> rabbitmqctl list_users
superuser	[administrator]
MLUser	[consumer]
~~~

Check declared queues:

~~~
>> rabbitmqctl list_queues
name	messages
response	1
manage	0
fib	1
~~~