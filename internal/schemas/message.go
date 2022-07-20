package schemas

import (
	"github.com/streadway/amqp"
)

// MessageRequest Message is RequestBody from Client
type MessageRequest struct {
	CorrelationId string
	Data          string `json:"data" binding:"required"`
	Headers       amqp.Table
}

type MessageCreate struct {
	//CorrelationId string
	//Body          []byte
	//Headers       amqp.Table
	//RoutingKey    string
	RmqMessage amqp.Delivery
}

// MessageReply is Response mapping queue
type MessageReply struct {
	CorrelationId string
	Data          []byte
	Headers       amqp.Table
}
