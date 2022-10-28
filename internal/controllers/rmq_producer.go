package controllers

import (
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
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
			log.Printf("PRODUCE: %s", conn.Queue)

			err := msg.MsgMq.Ack(false)
			if err != nil {
				log.Printf("ERROR: failed to ack message: %s", err.Error())
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
				log.Err(err).Msgf("ERROR: fail to publish msg: %s", msg.CorrelationId)
			}
			log.Printf("INFO: [%v] - published", msg.CorrelationId)

		case errCh := <-DLEChannel:

			if errCh.DLEStop != true {
				err := errCh.MsgMq.Nack(false, false)
				if err != nil {
					log.Error().Err(err).Msg("Failed Nack msg")
				}
			} else {
				log.Info().Msgf("LIMIT: number of attempts to send a message has reached the maximum.\nMessage removed from dle")
				err := errCh.MsgMq.Ack(false)
				if err != nil {
					log.Printf("ERROR: failed to ack message: %s", err.Error())
				}
			}
		}
	}
}
