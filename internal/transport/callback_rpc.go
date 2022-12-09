package transport

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	"provider_mq/internal/controllers"
	"provider_mq/internal/schemas"

	modelrpc "provider_mq/internal/schemas/model"
)

func (c *HostConfig) waitRpcReplyModel(msg amqp.Delivery) {

	newMsg := &schemas.MessageRpc{
		Data: string(msg.Body),
	}
	newData, err := json.Marshal(newMsg)
	if err != nil {
		log.Error().Err(err).Msg("Failed parse json body")
	}

	conn, err := grpc.Dial(
		c.host+":"+c.port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error().Err(err).Msgf("error dial to grpc channel")
		msgReply := &schemas.MessageReplyError{
			Error: err,
			MsgMq: msg,
		}
		controllers.DLEChannel <- *msgReply
		return
	}
	defer conn.Close()

	client := modelrpc.NewBaseModelInterfaceClient(conn)
	request := &modelrpc.ModelInterfaceRequest{
		Data: newData,
	}
	response, err := client.Predict(context.Background(), request)

	if err != nil {
		grpclog.Errorf("gRPC service unavailable", err)
		dleStop := canRetry(msg.Headers)
		msgReply := &schemas.MessageReplyError{
			Error:   err,
			MsgMq:   msg,
			DLEStop: dleStop,
		}
		controllers.DLEChannel <- *msgReply

	} else {

		headers := make(amqp.Table)

		msgReply := &schemas.MessageReplySuccess{
			CorrelationId: msg.CorrelationId,
			Data:          response.GetData(),
			Headers:       headers,
			MsgMq:         msg,
		}

		controllers.PublishChannel <- *msgReply
	}
}
