package controllers

import "github.com/rs/zerolog/log"

//delete(ReplyChannels, msgRequest.CorrelationId)

func ProduceMessage() {
	for {
		select {
		case msg := <-PublishChannel:
			log.Info().Msgf("Message to Publish: %v", msg)
		}
	}
}
