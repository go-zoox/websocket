package conn

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

type Conn interface {
	Close() error
	//
	WriteMessage(typ int, msg []byte) error
	//
	WriteTextMessage(msg []byte) error
	WriteBinaryMessage(msg []byte) error
	//
	Ping(msg []byte) error
	Pong(msg []byte) error
}

const (
	TextMessage   = websocket.TextMessage
	BinaryMessage = websocket.BinaryMessage
)

type conn struct {
	ctx context.Context
	raw *websocket.Conn
}

func New(ctx context.Context, raw *websocket.Conn) Conn {
	return &conn{
		ctx: ctx,
		raw: raw,
	}
}

func (c *conn) Context() context.Context {
	return c.ctx
}

func (c *conn) Close() error {
	return c.raw.Close()
}

func (c *conn) WriteMessage(typ int, msg []byte) error {
	return c.raw.WriteMessage(typ, msg)
}

func (c *conn) WriteTextMessage(msg []byte) error {
	return c.raw.WriteMessage(TextMessage, msg)
}

func (c *conn) WriteBinaryMessage(msg []byte) error {
	return c.raw.WriteMessage(BinaryMessage, msg)
}

func (c *conn) Ping(msg []byte) error {
	return c.raw.WriteControl(websocket.PingMessage, msg, time.Now().Add(time.Second))
}

func (c *conn) Pong(msg []byte) error {
	return c.raw.WriteControl(websocket.PongMessage, msg, time.Now().Add(time.Second))
}
