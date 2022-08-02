package service

import (
	"strings"

	"provider_mq/internal/app"
	"provider_mq/internal/consts"
	"provider_mq/internal/controllers"
	"provider_mq/internal/schemas"
	"provider_mq/internal/utils"
)

type HostConfig struct {
	host     string
	port     string
	basePath string
}

func Predict() {

	//errg := new(errgroup.Group)

	go func() {
		app.SetupApp()
	}()

	modelHost := utils.GetEnvVar("MODEL_HOST")
	modelPort := utils.GetEnvVar("MODEL_PORT")
	if modelHost == "" {
		checkModelHostEnv := utils.GetEnvVar("MODEL_HOST_ENV")
		if checkModelHostEnv != "" {
			modelHost = utils.GetEnvVar(strings.ToUpper(checkModelHostEnv))
		} else {
			modelHost = consts.DefaultHostModel
		}
	}
	if modelPort == "" {
		modelPort = consts.DefaultPortModel
	}

	hostModel := &HostConfig{
		host:     modelHost,
		port:     modelPort,
		basePath: consts.BasePath,
	}

	for {
		select {
		case msgRequest := <-controllers.RequestChannel:

			msg := msgRequest

			go func(msg schemas.MessageCreate) {
				hostModel.waitReplyModel(msg.RmqMessage)
			}(msg)
		}
	}
}
