package client

import (
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
)

func (c *client) OnMessage(fn func(conn conn.Conn, mt int, message []byte) error) {
	c.ee.On(event.TypeMessage, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadMessage)
		if !ok {
			c.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := fn(p.Conn, p.Type, p.Message); err != nil {
			c.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  p.Conn,
				Error: err,
			})
		}
	}))
}

func (c *client) OnTextMessage(cb func(conn conn.Conn, message []byte) error) {
	c.OnMessage(func(connx conn.Conn, mt int, message []byte) error {
		if mt == conn.TextMessage {
			return cb(connx, message)
		}

		return nil
	})
}

func (c *client) OnBinaryMessage(cb func(conn conn.Conn, message []byte) error) {
	c.OnMessage(func(connx conn.Conn, mt int, message []byte) error {
		if mt == conn.BinaryMessage {
			return cb(connx, message)
		}

		return nil
	})
}
