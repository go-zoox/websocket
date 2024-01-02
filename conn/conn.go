package conn

import (
	"context"
	"net/http"
	"time"

	"github.com/go-zoox/core-utils/safe"
	"github.com/go-zoox/eventemitter"
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
	ReadMessage() (int, []byte, error)
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
	//
	On(typ string, handler eventemitter.Handle)
	Emit(typ string, payload any)
	//
	Get(key string) any
	Set(key string, value any) error
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
	//
	ee    *eventemitter.EventEmitter
	cache *safe.Map
}

func New(ctx context.Context, raw *websocket.Conn, req *http.Request) Conn {
	return &conn{
		id:  uuid.V4(),
		ctx: ctx,
		raw: raw,
		req: req,
		//
		ee:    eventemitter.New(),
		cache: safe.NewMap(),
	}
}

func (c *conn) ID() string {
	return c.id
}

func (c *conn) Context() context.Context {
	return c.ctx
}

func (c *conn) Close() error {
	c.ee.Stop()
	return c.raw.Close()
}

func (c *conn) WriteMessage(typ int, msg []byte) error {
	return c.raw.WriteMessage(typ, msg)
}

func (c *conn) WriteTextMessage(msg []byte) error {
	return c.WriteMessage(TextMessage, msg)
}

func (c *conn) WriteBinaryMessage(msg []byte) error {
	return c.WriteMessage(BinaryMessage, msg)
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

func (c *conn) On(typ string, handler eventemitter.Handle) {
	c.ee.On(typ, handler)
}

func (c *conn) Emit(typ string, payload any) {
	c.ee.Emit(typ, payload)
}

func (c *conn) Get(key string) any {
	return c.cache.Get(key)
}

func (c *conn) Set(key string, value any) error {
	return c.cache.Set(key, value)
}

func (c *conn) ReadMessage() (int, []byte, error) {
	return c.raw.ReadMessage()
}
