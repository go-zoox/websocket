package client

import (
	"errors"
	"fmt"

	"github.com/go-zoox/logger"
	"github.com/go-zoox/safe"
	"github.com/go-zoox/uuid"
	"github.com/go-zoox/websocket"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/extension/event/entity"
)

type Client interface {
	Emit(event string, payload entity.EventPayload, callback func(err error, payload entity.EventPayload), opts ...EmitOption) error
}

type client struct {
	core websocket.Client
	//
	cbs map[string]EventCallback
}

type EventCallback struct {
	Callback    func(err error, payload entity.EventPayload)
	IsSubscribe bool
}

func New(core websocket.Client) Client {
	c := &client{
		core: core,
		cbs:  make(map[string]EventCallback),
	}

	// event
	core.OnTextMessage(func(conn conn.Conn, message []byte) error {
		go func() {
			event := &entity.Event{}
			if err := event.Decode(message); err != nil {
				logger.Errorf("[event] failed to decode: %s (message: %s)", err, string(message))
				// return err
				return
			}

			if fn, ok := c.cbs[event.ID]; ok {
				if !fn.IsSubscribe {
					delete(c.cbs, event.ID)
				}

				err := safe.Do(func() error {
					var err error
					if event.Error != "" {
						err = errors.New(event.Error)
					}

					fn.Callback(err, event.Payload)
					return nil
				})
				if err != nil {
					logger.Errorf("[event] failed to handle event callback: %v", err)
				}
			}
		}()

		return nil
	})

	return c
}

type EmitConfig struct {
	IsSubscribe bool
}

type EmitOption func(cfg *EmitConfig)

func (c *client) Emit(typ string, payload entity.EventPayload, callback func(err error, payload entity.EventPayload), opts ...EmitOption) error {
	cfg := &EmitConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	event := &entity.Event{
		ID:      uuid.V4(),
		Type:    typ,
		Payload: payload,
	}

	bytes, err := event.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode event: %s", err)
	}

	// subscribe event -> no need to wait for response
	if cfg.IsSubscribe {
		c.cbs[event.ID] = EventCallback{
			Callback:    callback,
			IsSubscribe: true,
		}

		return c.core.SendMessage(bytes)
	}

	// non-subscribe event -> wait for response
	done := make(chan struct{})
	c.cbs[event.ID] = EventCallback{
		Callback: func(err error, payload entity.EventPayload) {
			callback(err, payload)
			done <- struct{}{}
		},
	}

	if err := c.core.SendMessage(bytes); err != nil {
		return fmt.Errorf("failed to write text message: %s", err)
	}

	<-done
	return nil
}
