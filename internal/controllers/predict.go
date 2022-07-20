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
	//var msg schemas.MessageRequest

	//requestIdHeaderName := consts.RequestIdHttpHeaderName
	//requestId := c.GetString(requestIdHeaderName)

	// testing

	for {
		select {
		case msgRequest := <-RequestChannel:

			log.Printf("INFO: [%s] received", msgRequest.RmqMessage.CorrelationId)

			//for k := range msgRequest.Headers {
			//	if str, ok := msgRequest.Headers[k].(string); ok {
			//		router.Header().Set(k, str)
			//	}
			//}

			//Encode the data
			//postBody, _ := json.Marshal(map[string]string{
			//	"name":  "Toby",
			//	"email": "Toby@example.com",
			//})
			//msgToModel := bytes.NewBuffer(postBody)

			waitReplyModel(msgRequest.RmqMessage.Body, msgRequest.RmqMessage) // answer =

			// PublishChannel <- *msgRequest
			// add case <- channel to Errors form http

		}
	}
}

func waitReplyModel(msgToModel []byte, msg amqp.Delivery) {

	req, err := http.NewRequest(
		"POST",
		fmt.Sprint("http://"+consts.HostModel+":"+consts.PortModel+consts.BasePath),
		bytes.NewBuffer(msgToModel),
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

	//msg.Ack(true)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		log.Error().Err(err).Msg("Timeout on send message to <Model Application>")
	} else if err != nil {
		log.Error().Err(err).Msg("Error on send message to <Model Application>")
		msg.Nack(true, false)
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
}
