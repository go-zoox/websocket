package event

import (
	"fmt"

	"github.com/go-zoox/websocket/conn"
)

type PayloadError struct {
	Conn  conn.Conn
	Error error
}

const TypeError = "[internal] error"

var ErrInvalidPayload = &PayloadError{
	Error: fmt.Errorf("invalid payload"),
}
