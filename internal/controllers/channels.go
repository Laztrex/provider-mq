package controllers

import (
	"provider_mq/internal/schemas"
)

// RequestChannel PublishChannels channel to publish rabbit messages
var RequestChannel = make(chan schemas.MessageCreate, 10)

var PublishChannel = make(chan schemas.MessageReplySuccess, 10)

var DLEChannel = make(chan schemas.MessageReplyError, 10)
