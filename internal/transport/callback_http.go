package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"io"
	"net"
	"net/http"
	"time"

	"provider_mq/internal/consts"
	"provider_mq/internal/controllers"
	"provider_mq/internal/middlewares"
	"provider_mq/internal/schemas"
)

func (c *HostConfig) WaitReplyModel(msg amqp.Delivery) {

	newMsg := make(schemas.MessageRest)
	err := json.Unmarshal(msg.Body, &newMsg)

	newData, err := json.Marshal(newMsg)
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

	if err != nil {

		pushFailureResponse(&msg, resp, err)

	} else {

		pushSuccessResponse(&msg, resp)

	}
}

func pushSuccessResponse(msg *amqp.Delivery, resp *http.Response) {
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
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
		MsgMq:         *msg,
		RequestId:     resp.Header.Get(consts.KeyRequestId),
	}

	controllers.PublishChannel <- *msgReply
}

func pushFailureResponse(msg *amqp.Delivery, resp *http.Response, err error) {
	dleStop := canRetry(msg.Headers)
	msgReply := &schemas.MessageReplyError{
		Error:     err,
		MsgMq:     *msg,
		DLEStop:   dleStop,
		RequestId: resp.Header.Get(consts.KeyRequestId),
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		log.Error().Err(err).Msgf("Timeout on send request to <Model Application>: %s", err.Error())
	} else {
		log.Error().Err(err).Msgf("Error while send request to <Model Application>: %s", err.Error())
	}

	controllers.DLEChannel <- *msgReply
}
