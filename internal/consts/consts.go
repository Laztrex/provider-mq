package consts

const EnvFile = ".env"
const EnvFileDirectory = "."

const QueuesConf = "configs/queue_config.yaml"

const DefaultHostModel = "localhost"
const DefaultPortModel = "8080"
const BasePath = "/predict"
const RequestIdHttpHeaderName = "request-id"

const RestTimeout = 300
const DleRetry = 5

const LogPath = "/var/log/metrics.log"

const (
	MqCACERT = "/certs/provider_mq/cacert.pem"
	MqCERT   = "/certs/provider_mq/client_cert.pem"
	MqKEY    = "/certs/provider_mq/client_key.pem"
)
