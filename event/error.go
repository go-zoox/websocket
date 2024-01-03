package event

import (
	"fmt"
)

type PayloadError struct {
	Error error
}

const TypeError = "[internal] error"

var ErrInvalidPayload = &PayloadError{
	Error: fmt.Errorf("invalid payload"),
}
