package client

import (
	"fmt"

	"github.com/go-zoox/uuid"
	"github.com/go-zoox/websocket/event/cs"
)

// Event triggers a event with handler
func (c *client) Event(typ string, payload any, callback func(*cs.EventPayload)) error {
	event := &cs.Event{
		ID:   uuid.V4(),
		Type: typ,
		Payload: cs.EventPayload{
			Data: payload,
		},
	}

	bytes, err := event.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode event: %s", err)
	}

	c.cbs.events[event.ID] = callback

	return c.conn.WriteTextMessage(bytes)
}
