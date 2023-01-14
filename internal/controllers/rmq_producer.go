package controllers

import (
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	"provider_mq/internal/consts"
)

func (conn *RMQSpec) ProduceDeclare() {

	err := conn.QueueDeclare()
	conn.OnError(err, "Failed to declare a queue while consuming")

}

func (conn *RMQSpec) ProduceMessages() {
	for {
		select {
		case err := <-conn.Err:
			err = conn.Reconnect()
			if err != nil {
				panic(err)
			}

		case msg := <-PublishChannel:
			//var rcvRequestId string
			//if _, ok := msg.Headers[consts.KeyRequestId]; ok {
			//	rcvRequestId = msg.Headers[consts.KeyRequestId].(string)
			//} else {
			//	rcvRequestId = ""
			//}

			log.Info().
				Str(consts.KeyRequestId, msg.RequestId).
				Msgf("PRODUCE: %s", conn.Queue)

			err := msg.MsgMq.Ack(false)
			if err != nil {
				log.Error().Err(err).
					Str(consts.KeyRequestId, msg.RequestId).
					Msgf("ERROR: failed to ack message: %s", err.Error())
			}

			err = conn.Channel.Publish(
				"", // exchange
				//conn.RoutingKey, // routing key
				msg.MsgMq.ReplyTo,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					Body:          msg.Data,
					Headers:       msg.Headers,
					CorrelationId: msg.CorrelationId,
				},
			)
			if err != nil {
				log.Error().Err(err).
					Str(consts.KeyCorrelationId, msg.CorrelationId).
					Msg("ERROR: fail to publish msg")
			}
			log.Info().
				Str(consts.KeyCorrelationId, msg.CorrelationId).
				Msg("message published")

		case errMsg := <-DLEChannel:

			if errMsg.DLEStop != true {
				err := errMsg.MsgMq.Nack(false, false)
				if err != nil {
					log.Error().Err(err).
						Str(consts.KeyCorrelationId, errMsg.MsgMq.CorrelationId).
						Msg("Failed Nack msg")
				}
			} else {
				log.Info().
					Str(consts.KeyCorrelationId, errMsg.MsgMq.CorrelationId).
					Msgf("LIMIT: number of attempts to send a message has reached the maximum.\n" +
						"Message removed from dle")
				err := errMsg.MsgMq.Ack(false)
				if err != nil {
					log.Error().Err(err).
						Str(consts.KeyCorrelationId, errMsg.MsgMq.CorrelationId).
						Msgf("ERROR: failed to ack message: %s", err.Error())
				}
			}
		}
	}
}
