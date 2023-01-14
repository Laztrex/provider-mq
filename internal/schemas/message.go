package schemas

import (
	"github.com/streadway/amqp"
)

type MessageRest map[string]interface{}

type MessageCreate struct {
	RmqMessage amqp.Delivery
	RequestId  string
}

type MessageRpc struct {
	Data string `json:"data"`
}

// MessageReplySuccess is Response mapping queue
type MessageReplySuccess struct {
	MsgMq         amqp.Delivery
	CorrelationId string
	Data          []byte
	Headers       amqp.Table
	RequestId     string
}

// MessageReplyError is Response if error request to Model
type MessageReplyError struct {
	MsgMq     amqp.Delivery
	Error     error
	DLEStop   bool
	RequestId string
}
