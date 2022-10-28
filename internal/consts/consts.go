package consts

const (
	EnvFile          = ".env"
	EnvFileDirectory = "."

	QueuesConf = "configs/queue_config.yaml"

	DefaultHostModel = "webapp"
	DefaultPortModel = "8080"

	BasePath                = "/query"
	RequestIdHttpHeaderName = "request-id"

	RestTimeout = 300
	DleRetry    = 5

	LogPath = "/var/log/provider-mq-metrics.log"

	MqCACERT = "/certs/provider_mq/cacert.pem"
	MqCERT   = "/certs/provider_mq/client_cert.pem"
	MqKEY    = "/certs/provider_mq/client_key.pem"
)
