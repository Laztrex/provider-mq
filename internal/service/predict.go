package service

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net"
	"net/http"
	"provider_mq/internal/consts"
	"provider_mq/internal/controllers"
	"provider_mq/internal/schemas"
	"time"
)

func Predict() {

	//errg := new(errgroup.Group)

	for {
		select {
		case msgRequest := <-controllers.RequestChannel:

			log.Printf("INFO: [%s] received", msgRequest.RmqMessage.CorrelationId)

			msg := msgRequest

			go func(msg schemas.MessageCreate) {
				waitReplyModel(msg.RmqMessage)
			}(msg)
		}
	}
}

func waitReplyModel(msg amqp.Delivery) {

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

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		log.Error().Err(err).Msgf("Timeout on send message to <Model Application>: %s", err.Error())
		msgReply := &schemas.MessageReplyError{
			Error: err,
			MsgMq: msg,
		}
		controllers.DLEChannel <- *msgReply
	} else {

		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error read Response")
		}

		headers := make(amqp.Table)

		for k, v := range resp.Header {
			headers[k] = v[0]
			//headers[textproto.CanonicalMIMEHeaderKey(k)] = v[0]
		}

		msgReply := &schemas.MessageReplySuccess{
			CorrelationId: msg.CorrelationId,
			Data:          response,
			Headers:       headers,
			MsgMq:         msg,
		}

		controllers.PublishChannel <- *msgReply
	}
}
