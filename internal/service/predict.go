package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net"
	"net/http"
	"provider_mq/internal/app"
	"provider_mq/internal/consts"
	"provider_mq/internal/controllers"
	"provider_mq/internal/middlewares"
	"provider_mq/internal/schemas"
	"provider_mq/internal/utils"
	"time"
)

type HostConfig struct {
	host     string
	port     string
	basePath string
}

func Predict() {

	//errg := new(errgroup.Group)

	go func() {
		app.SetupApp()
	}()

	modelHost := utils.GetEnvVar("MODEL_HOST")
	modelPort := utils.GetEnvVar("MODEL_PORT")
	if modelHost == "" {
		modelHost = consts.HostModel
	}
	if modelPort == "" {
		modelPort = consts.PortModel
	}
	hostModel := &HostConfig{
		host:     modelHost,
		port:     modelPort,
		basePath: consts.BasePath,
	}

	for {
		select {
		case msgRequest := <-controllers.RequestChannel:

			log.Printf("INFO: [%s] received", msgRequest.RmqMessage.CorrelationId)

			msg := msgRequest

			go func(msg schemas.MessageCreate) {
				hostModel.waitReplyModel(msg.RmqMessage)
			}(msg)
		}
	}
}

func (c *HostConfig) waitReplyModel(msg amqp.Delivery) {

	newMsg := &schemas.MessageRest{
		Data: string(msg.Body),
	}

	newData, err := json.Marshal(*newMsg)
	if err != nil {
		log.Error().Err(err).Msg("Failed parse json body")
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprint("http://"+c.host+":"+c.port+c.basePath),
		bytes.NewBuffer(newData),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed create request")
	}

	for k := range msg.Headers {
		if str, ok := msg.Headers[k].(string); ok {
			req.Header.Set(k, str)
		}
	}

	client := &http.Client{
		Timeout:   consts.RestTimeout * time.Second,
		Transport: middlewares.LoggerTransportRoundTripper{Proxy: http.DefaultTransport},
	}

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
