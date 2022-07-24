package schemas

import (
	"github.com/streadway/amqp"
)

type MessageCreate struct {
	RmqMessage amqp.Delivery
}

// MessageReplySuccess is Response mapping queue
type MessageReplySuccess struct {
	MsgMq         amqp.Delivery
	CorrelationId string
	Data          []byte
	Headers       amqp.Table
}

// MessageReplyError is Response if error request to Model
type MessageReplyError struct {
	MsgMq  amqp.Delivery
	Error  error
	DLECnt int
}
