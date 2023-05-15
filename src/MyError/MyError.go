package MyError

import (
	"errors"
	"fmt"
)

var (
	ErrUndefinedMessageType = func(messageType string) error {
		return fmt.Errorf("undefined message type:%s", messageType)
	}
	ErrConnectionClosed = errors.New("the connection has not been built or has already been closed")
	ErrReconnectFailed  = errors.New("reconnect failed")
	ErrInvalidRequsetType = errors.New("invalid request type")
)
