package conn

import (
	"context"
	"net/http"
	"time"

	"github.com/go-zoox/uuid"
	"github.com/gorilla/websocket"
)

type Conn interface {
	ID() string
	//
	Context() context.Context
	//
	Close() error
	//
	WriteMessage(typ int, msg []byte) error
	//
	WriteTextMessage(msg []byte) error
	WriteBinaryMessage(msg []byte) error
	//
	Ping(msg []byte) error
	Pong(msg []byte) error
	//
	Raw() *websocket.Conn
	//
	Request() *http.Request
}

const (
	TextMessage   = websocket.TextMessage
	BinaryMessage = websocket.BinaryMessage
)

type conn struct {
	id  string
	ctx context.Context
	raw *websocket.Conn
	req *http.Request
}

func New(ctx context.Context, raw *websocket.Conn, req *http.Request) Conn {
	return &conn{
		id:  uuid.V4(),
		ctx: ctx,
		raw: raw,
		req: req,
	}
}

func (c *conn) ID() string {
	return c.id
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

func (c *conn) Raw() *websocket.Conn {
	return c.raw
}

func (c *conn) Request() *http.Request {
	return c.req
}
