package service

import (
	"provider_mq/internal/app"
	"provider_mq/internal/controllers"
	"provider_mq/internal/schemas"
)

func Predict() {

	go app.SetupApp()

	hostModel := getHostModel()

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
