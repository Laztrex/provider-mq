package transport

import (
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"strings"

	"provider_mq/internal/consts"
	"provider_mq/internal/utils"
)

type HostConfig struct {
	host     string
	port     string
	basePath string
}

func GetHostModel() *HostConfig {

	modelHost := utils.GetEnvVar("MODEL_HOST")
	modelPort := utils.GetEnvVar("MODEL_PORT")

	if modelHost == "" {
		checkModelHostEnv := utils.GetEnvVar("MODEL_HOST_ENV")
		if checkModelHostEnv != "" {
			modelHost = utils.GetEnvVar(strings.ToUpper(checkModelHostEnv))
		} else {
			modelHost = consts.HostModelDefault
		}
	}
	if modelPort == "" {
		modelPort = consts.PortModelDefault
	}

	hostModel := &HostConfig{
		host:     modelHost,
		port:     modelPort,
		basePath: consts.BasePath,
	}

	return hostModel

}

// canRetry checks the number of attempts to release dle
func canRetry(h amqp.Table) bool {

	if vDle, ok := h["x-death"]; ok {
		if j, ok := vDle.([]interface{}); ok {
			if dleCount, ok := j[0].(amqp.Table)["count"]; ok {
				log.Info().Msgf("value count: %v, type count: %T", dleCount, dleCount)
				if dleCount.(int64) > consts.DleRetry {
					return true
				}
			}
		}
	}
	return false
}
