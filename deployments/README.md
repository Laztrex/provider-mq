# Provider-gateway-mq Sample on Openshift

Here is an example of deploying a controller in Openshift.  
The steps in this document assume that you have access to an OpenShift deployment to which you can deploy applications, including a default set of image streams. Instructions for setting the default image streams are available here. If you are now defining a set of image streams, be sure to pass in the appropriate cluster administrator credentials and create the image streams in the "OpenShift" namespace.

## Description

The project is part of a cloud-based microservice solution for interacting with an ML-application and related infrastructure components.  
The *Provider-gateway-mq* controller is the link for the interaction between the Initiator and the ML model.

![Image alt](https://github.com/Laztrex/provider-mq/blob/master/docs/pics/sketch_provider_mq.png)

Naturally, for more reliable and stable operation, it is necessary to provide Envoy Proxy connection, as well as Service Mesh concept services for traffic control and monitoring.  
This example is just a demonstration.

You can also install [provider-**gateway**-mq](https://github.com/Laztrex/provider-gateway-mq) that allows you to make REST-requests:

![Image alt](https://github.com/Laztrex/provider-mq/blob/master/docs/pics/with_provider_gateway_mq.png)

### Functionality

* provider-mq should be responsible for establishing a secure connection between the message broker server (directly) and/or Initiator (indirectly) and model instances
* provider-mq should be responsible for establishing a secure connection between the Initiator and auxiliary services
* The possibility of declaring queues/topics in provider-mq, which will serve for publishing and consuming messages
* Interaction between:
  * Initiator and provider-mq occur asynchronously - data is directly published to the external server (or StatefulSets) of the message broker
  * provider-mq and ML model, if an event is initiated by the Initiator, it occurs in synchronous mode - for each message, the model must send the result of the operation in synchronous mode
  * Confirmation of message processing occurs due to the presence of two multidirectional queues queue.in/queue.out. In addition to pub/sub interaction, the organization of an RPC pattern is provided
* Reading messages and executing models is concurrent
* In the case when the model worked incorrectly and ended with an error, the result will be published to the output channel of the message broker, in addition, the metric about the unsuccessful execution of the model is discarded
  * If the *dle* parameter is specified in the configuration file for the controller-mq, in case of a timeout (300s by default) or an error, the message is discarded to a dedicated queue with a pending call
  

## Deployment steps

* Set roles and access rights (taken from the [example](https://kubernetes.io/docs/reference/access-authn-authz/rbac/))
  ~~~yaml
  ---
  apiVersion: v1
  kind: ServiceAccount
  metadata:
    name: rabbitmq
  ---
  kind: Role
  apiVersion: rbac.authorization.k8s.io/v1beta1
  metadata:
  name: endpoint-reader
  rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get"]
  ---
  kind: RoleBinding
  apiVersion: rbac.authorization.k8s.io/v1beta1
  metadata:
    name: endpoint-reader
  subjects:
  - kind: ServiceAccount
    name: rabbitmq
  roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: Role
    name: endpoint-reader
  ~~~
  
* Persistent storage
  If the RabbitMQ configuration is set in the StatefulSet in the project, then it is worth providing persistent storage
  
* Services

  - Create a *headless*-service for the Peer Discovery plugin (for StatefulSets application)
    ~~~
    oc create -f conf/service-headless.yml
    ~~~
  - Create a Service for external access to Gateway-mq
    ~~~
    oc create -f conf/service-provider-mq.yml
    ~~~
  - Service for external applications with RabbitMQ, incl. using the admin panel
    ~~~
    oc create -f conf/service-rabbitmq.yml
    ~~~
  - Service for ML-application
    ~~~
    oc create -f conf/service-ml.yml
    ~~~

* Configuration
  
  - Creating RabbitMQ Configuration Files
    ~~~
    oc create -f conf/configmap.yml
    ~~~
    
  - Create a configuration file for *Provider-gateway-mq*
    ~~~
    oc create -f conf/configmap-mq.yml
    ~~~
  
* Applications
 
  - Set the StatefulSet configuration of the RabbitMQ client
    ~~~
    oc create -f conf/statefulset.yml
    ~~~

  - Deploying the *Gateway-MQ* Application
    ~~~
    oc process -f conf/dc-gateway.yml --ignore-unknown-parameters --param-file="templates/template_mlx.env" | oc apply -f -
    ~~~

  - Deploying the *Provider-MQ* Application
    ~~~
    oc process -f conf/dc-provider.yml --ignore-unknown-parameters --param-file="templates/template_mlx.env" | oc apply -f -
    ~~~

  - Deploying ML Application
    ~~~
    oc process -f conf/dc-mlx.yml --ignore-unknown-parameters --param-file="templates/template_mlx.env" | oc apply -f -
    ~~~

![Image alt](https://github.com/Laztrex/provider-mq/blob/master/docs/pics/get_pods.png)

> Fluentd sidecar is omitted from these sample config files.

Checking:
Publish message in Rabbitmq cluster.
Or used *[provider-gateway-mq](https://github.com/Laztrex/provider-gateway-mq)*

~~~
curl -k --key client_key.pem --cert client_cert.pem -d '{"data": "[29, 29]"}' -H "RqUID: 52-52-52-52" -H "Content-Type: application/json" -H "routing-key: predict.now" https://mlx.my-fqdn/v1/predict
{"result": 42}
~~~
