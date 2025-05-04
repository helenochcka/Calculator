package rabbitmq

import "errors"

var ErrUnmarshallingMsg = errors.New("failed to unmarshal message")
var ErrConsumingMsgs = errors.New("failed to consume messages")
var ErrCalculatingResult = errors.New("")
