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

			//r := reflect.ValueOf(msg.Headers)
			//requestId = reflect.Indirect(r).FieldByName(consts.RequestIdHttpHeaderName).String()
			//if requestId == "" {
			//	requestId = uuid.New().String()
			//}

			msgRequest := &schemas.MessageCreate{
				RmqMessage: msg,
			}

			//replyChannel := make(chan schemas.MessageRequest, 10)
			//ReplyChannels[msg.CorrelationId] = replyChannel

			//err = msg.Ack(true)
			//if err != nil {
			//	log.Printf("ERROR: fail to ack: %s", err.Error())
			//}

			RequestChannel <- *msgRequest

		}
	}
}
