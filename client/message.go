package client

import (
	"github.com/go-zoox/websocket/conn"
)

func (c *client) OnMessage(fn func(conn conn.Conn, mt int, message []byte) error) {
	c.cbs.messages = append(c.cbs.messages, fn)
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

func (c *client) SendMessage(typ int, message []byte) error {
	return c.conn.WriteMessage(typ, message)
}

func (c *client) SendTextMessage(message []byte) error {
	return c.conn.WriteMessage(conn.TextMessage, message)
}

func (c *client) SendBinaryMessage(message []byte) error {
	return c.conn.WriteMessage(conn.BinaryMessage, message)
}
