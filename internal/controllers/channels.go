package controllers

import (
	"provider_mq/internal/schemas"
)

// ReplyChannels RequestChannel keep waiting channels for reply messages from rabbit
var ReplyChannels = make(map[string]chan schemas.MessageRequest)

// RequestChannel PublishChannels channel to publish rabbit messages
var RequestChannel = make(chan schemas.MessageCreate, 10)

var PublishChannel = make(chan schemas.MessageReply, 10)
