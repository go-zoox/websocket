package client

import (
	"fmt"

	"github.com/go-zoox/uuid"
	"github.com/go-zoox/websocket/event/cs"
)

type EventConfig struct {
	IsSubscribe bool
}

type EventOption func(cfg *EventConfig)

// Event triggers a event with handler
func (c *client) Event(typ string, payload cs.EventPayload, callback func(err error, payload cs.EventPayload), opts ...EventOption) error {
	cfg := &EventConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	event := &cs.Event{
		ID:      uuid.V4(),
		Type:    typ,
		Payload: payload,
	}

	bytes, err := event.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode event: %s", err)
	}

	c.cbs.events[event.ID] = EventCallback{
		Callback:    callback,
		IsSubscribe: cfg.IsSubscribe,
	}

	return c.conn.WriteTextMessage(bytes)
}
