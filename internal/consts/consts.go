package consts

const (
	EnvFile          = ".env"
	EnvFileDirectory = "."

	QueuesConf = "configs/queue_config.yaml"

	HostModelDefault = "webapp"
	PortModelDefault = "8080"

	BasePath = "/query"

	RestTimeout = 300
	DleRetry    = 5

	MqCaCertDefault  = "/certs/provider_mq/cacert.pem"
	MqCertDefault    = "/certs/provider_mq/client_cert.pem"
	MqCertKeyDefault = "/certs/provider_mq/client_key.pem"

	LogPath = "/var/log/metrics.log"
)

const (
	KeyRequestId     = "RqUID"
	KeyCorrelationId = "CorrID"
)
