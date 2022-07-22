package controllers

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net"
	"net/http"
	"provider_mq/internal/consts"
	"provider_mq/internal/schemas"
)

func Predict() {

	for {
		select {
		case msgRequest := <-RequestChannel:

			log.Printf("INFO: [%s] received", msgRequest.RmqMessage.CorrelationId)

			go func(msg schemas.MessageCreate) {
				err := waitReplyModel(msg.RmqMessage)
				if err != nil {
					err := msg.RmqMessage.Nack(false, false)
					if err != nil {
						log.Error().Err(err).Msg("Failed Nack msg")
					}
				}
			}(msgRequest)

			// PublishChannel <- *msgRequest
			// add case <- channel to Errors form http

		}
	}
}

func waitReplyModel(msg amqp.Delivery) error {

	req, err := http.NewRequest(
		"POST",
		fmt.Sprint("http://"+consts.HostModel+":"+consts.PortModel+consts.BasePath),
		bytes.NewBuffer(msg.Body),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed create request")
	}

	for k := range msg.Headers {
		if str, ok := msg.Headers[k].(string); ok {
			req.Header.Set(k, str)
		}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	fmt.Print(resp)
	defer resp.Body.Close()

	if err, ok := err.(net.Error); ok && err.Timeout() {
		log.Info().Msg("Timeout on send message to <Model Application>")
		//err := msg.Nack(false, false)

		if err != nil {
			return fmt.Errorf("timeout on send message to <Model Application> %s", err)
		}
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error read Response")
	}

	//delete(ReplyChannels, msg.CorrelationId)
	//msg.Ack(true)

	err = msg.Ack(true)
	if err != nil {
		log.Printf("ERROR: fail to ack: %s", err.Error())
	}

	headers := make(amqp.Table)

	for k, v := range resp.Header {
		headers[k] = v[0]
		//headers[textproto.CanonicalMIMEHeaderKey(k)] = v[0]
	}

	msgReply := &schemas.MessageReply{
		CorrelationId: msg.CorrelationId,
		Data:          response,
		Headers:       headers,
	}

	PublishChannel <- *msgReply

	return nil
}
