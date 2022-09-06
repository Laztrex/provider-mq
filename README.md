# Provider-MQ

The project is part of a cloud-based microservice solution for interacting with an ML-application and related
infrastructure components.  


This application allows you to call for execution an ML model located in a dedicated server.   
*Consumer* is inside the application and when a message appears, a REST-request is made in the application with the ML model.

## Description

*Prover-MQ* currently implements the following:

* The initiator publishes a message
* After publishing - REST request to *ML model* service
  * *provider-mq* establishes a TLS connection to the message broker server (RabbitMQ)
* Application arranges sync/async calls with the model
* Defining queue parameters in a configuration file
* *ML model* starts the calculation, the result is returned synchronously to *provider-mq*
  * If the queue for responses (queue.out*) is not set, it is assumed that the initiator does not expect a response from the model. In this case, the model is a worker
* *provider-mq* publishes the response to queue.out, it is assumed that the initiator listens to the output queue
* In case of an unsuccessful attempt to launch the model for execution, the response is sent to the waiting queue with a deferred call (queue.dle*)
* Declare of exchange, queues on the *provider-mq* side - we do not trust the model  

*queue names are defined abstractly.

The project will be updated.

Tested on
~~~
>> go version
go version go1.17.2 darwin/amd64
~~~

## Usage

The project uses global *environment variables*.

| Env Name       |                                Goal                                 |         Expected value example         |
|----------------|:-------------------------------------------------------------------:|:--------------------------------------:|
| RMQ_URL        |                      Host to connect RabbitMQ                       | "amqps://user:password@rabbitmq:5671/" |
| LOG_LEVEL      |                            Logging level                            |                "debug"                 |
| MODEL_HOST     |                      Host request to the model                      |               "0.0.0.0"                |
| MODEL_HOST_ENV | Defining the project host in namespace (used instead of MODEL_HOST) |         "MLX_MLX_SERVICE_HOST"         |
| MODEL_PORT     |                      Port request to the model                      |                 "8080"                 |

and *constants*.

| Const Name         |                    Goal                    |         Expected value example          |
|--------------|:------------------------------------------:|:---------------------------------------:|
| MqCACERT    |                     CA                     |   "/certs/provider_mq/cacert_mq.pem"    |
| MqCERT      |                Certificate                 | "/certs/provider_mq/client_cert_mq.pem" |
| MqKEY       |              Key certificate               | "/certs/provider_mq/client_key_mq.pem"  |
| EnvFile                 |         Local env file in project          |                 ".env"                  |
| EnvFileDirectory        |                Dir env file                |                   "."                   |
| QueuesConf              |     Path to configration Queue declare     |       "configs/queue_config.yaml"       |
| RequestIdHttpHeaderName |        Name Header for *request-id*        |                 ".env"                  |
| LogPath                 |             Path to dump logs              |         "/var/log/metrics.log"          |
| DefaultHostModel            |  Default url for the request to the model  |                "0.0.0.0"                |
| DefaultPortModel        | Default port for the request to the model  |                 "8080"                  |
| BasePath              | Base endpoint for the request to the model |               "/predict"                |
| RestTimeout | Holding time of the model response channel |                   300                   |

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
~~~

Currently, only Topic can be defined from the configuration file for Exchange (but Direct can also be defined). Other exchanges - will be supplemented.
Flexible settings for the queue - lifetime, autodelete, types, arguments and etc. currently not included in the configuration file, the parameters can be configured inside the code optionally.  

The directory [examples/webapp](https://github.com/Laztrex/provider-mq/blob/master/examples/webapp/) contains a simple web application for testing the project.

~~~
>> go version
go version go1.17.2 darwin/amd64
~~~
~~~
>> rabbitmqctl version
3.10.6
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

Check:

~~~
>> curl http://127.0.0.1:5080/health
{"status":"OK"}
~~~
Now publish message to rabbitmq server.

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