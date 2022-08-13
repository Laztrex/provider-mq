package service

import (
	"provider_mq/internal/app"
	"provider_mq/internal/controllers"
	"provider_mq/internal/schemas"
)

func Predict() {

	//errg := new(errgroup.Group)

	go func() {
		app.SetupApp()
	}()

	hostModel := getHostModel()

	for {
		select {
		case msgRequest := <-controllers.RequestChannel:

			msg := msgRequest

			go func(msg schemas.MessageCreate) {
				hostModel.waitRpcReplyModel(msg.RmqMessage)
			}(msg)
		}
	}
}
