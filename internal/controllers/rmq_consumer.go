package controllers

import (
	"github.com/rs/zerolog/log"
	"provider_mq/internal/schemas"
)

func (conn *RMQSpec) ConsumeDeclare() {
	var err error

	if conn.Exchange != "" {
		err = conn.ExchangeDeclare()
		conn.OnError(err, "Failed to declare exchange while publishing")
	}

	err = conn.QueueDeclare()
	conn.OnError(err, "Failed to declare a queue while publishing")

	if conn.Exchange != "" {
		err = conn.QueueBind()
		conn.OnError(err, "Failed to bind a queue while publishing")
	}
}

func (conn *RMQSpec) ConsumeMessages() {

	msgChannel, err := conn.Channel.Consume(
		conn.Queue, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	conn.OnError(err, "ERROR: fail create channel")

	for {
		select {
		case err := <-conn.Err:
			err = conn.Reconnect()
			if err != nil {
				panic(err)
			}

		case msg := <-msgChannel:

			if msg.CorrelationId == "" {
				continue // utils.GetCorrelationId(),
			}
			log.Info().Msgf("Receive body: %v", msg.Body)

			msgRequest := &schemas.MessageCreate{
				RmqMessage: msg,
			}

			RequestChannel <- *msgRequest

		}
	}
}
