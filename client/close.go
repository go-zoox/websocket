package client

import (
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
)

func (c *client) OnClose(cb func(conn conn.Conn) error) {
	c.ee.On(event.TypeClose, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadClose)
		if !ok {
			c.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := cb(p.Conn); err != nil {
			c.ee.Emit(event.TypeError, err)
		}
	}))
}

func (c *client) Close() error {
	return c.conn.Close()
}
