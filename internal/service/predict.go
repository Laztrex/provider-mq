package service

import (
	"provider_mq/internal/app"
	"provider_mq/internal/controllers"
	"provider_mq/internal/schemas"
	"provider_mq/internal/transport"
)

func Predict() {

	go app.SetupApp()

	hostModel := transport.GetHostModel()

	for {
		select {
		case msgRequest := <-controllers.RequestChannel:

			msg := msgRequest

			go func(msg schemas.MessageCreate) {
				hostModel.WaitReplyModel(msg.RmqMessage)
			}(msg)
		}
	}
}
