package rabbitmq

import "errors"

var FailedUnmarshalMsg = errors.New("failed to unmarshal message")
var FailedConsumeMsgs = errors.New("failed to consume messages")
var UnsuccessfulResult = errors.New("")
