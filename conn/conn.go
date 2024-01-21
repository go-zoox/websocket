package conn

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-zoox/core-utils/safe"
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/uuid"
	"github.com/go-zoox/websocket/event"
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
	On(typ string, handler eventemitter.Handler)
	Emit(typ string, payload any)
	//
	Get(key string) any
	Set(key string, value any) error
	//
	OnError(func(err error) error)
	OnClose(func(code int, message string) error)
	OnPing(func(appData []byte) error)
	OnPong(func(appData []byte) error)
	OnMessage(func(typ int, message []byte) error)
	//
	OnConnect(cb func() error)
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
	ee    eventemitter.EventEmitter
	cache *safe.Map
	//
	sync.Mutex
}

func New(ctx context.Context, raw *websocket.Conn, req *http.Request) Conn {
	ee := eventemitter.New(func(opt *eventemitter.Option) {
		opt.Context = ctx
	})

	return &conn{
		id:  uuid.V4(),
		ctx: ctx,
		raw: raw,
		req: req,
		//
		ee:    ee,
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
	return c.raw.Close()
}

func (c *conn) WriteMessage(typ int, msg []byte) error {
	c.Lock()
	defer c.Unlock()

	return c.raw.WriteMessage(typ, msg)
}

func (c *conn) WriteTextMessage(msg []byte) error {
	return c.WriteMessage(TextMessage, msg)
}

func (c *conn) WriteBinaryMessage(msg []byte) error {
	return c.WriteMessage(BinaryMessage, msg)
}

func (c *conn) Ping(msg []byte) error {
	c.Lock()
	defer c.Unlock()

	return c.raw.WriteControl(websocket.PingMessage, msg, time.Now().Add(time.Second))
}

func (c *conn) Pong(msg []byte) error {
	c.Lock()
	defer c.Unlock()

	return c.raw.WriteControl(websocket.PongMessage, msg, time.Now().Add(time.Second))
}

func (c *conn) Raw() *websocket.Conn {
	return c.raw
}

func (c *conn) Request() *http.Request {
	return c.req
}

func (c *conn) On(typ string, handler eventemitter.Handler) {
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

func (c *conn) OnError(cb func(err error) error) {
	c.On(event.TypeError, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadError)
		if !ok {
			// s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			logger.Errorf("invalid payload: %v", payload)
			return
		}

		if err := cb(p.Error); err != nil {
			logger.Errorf("failed to handle error: %v", err)
		}
	}))
}

func (c *conn) OnClose(cb func(code int, message string) error) {
	c.On(event.TypeClose, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadClose)
		if !ok {
			c.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := cb(p.Code, p.Message); err != nil {
			c.Emit(event.TypeError, &event.PayloadError{
				Error: err,
			})
		}
	}))
}

func (c *conn) OnPing(cb func(appData []byte) error) {
	c.On(event.TypePing, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadPing)
		if !ok {
			c.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := cb(p.Message); err != nil {
			c.Emit(event.TypeError, &event.PayloadError{
				Error: err,
			})
		}
	}))
}

func (c *conn) OnPong(cb func(appData []byte) error) {
	c.On(event.TypePong, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadPong)
		if !ok {
			c.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := cb(p.Message); err != nil {
			c.Emit(event.TypeError, &event.PayloadError{
				Error: err,
			})
		}
	}))
}

func (c *conn) OnMessage(cb func(typ int, message []byte) error) {
	c.On(event.TypeMessage, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadMessage)
		if !ok {
			c.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := cb(p.Type, p.Message); err != nil {
			c.Emit(event.TypeError, &event.PayloadError{
				Error: err,
			})
		}
	}))
}

func (c *conn) OnTextMessage(cb func(message []byte) error) {
	c.OnMessage(func(typ int, message []byte) error {
		if typ == TextMessage {
			return cb(message)
		}

		return nil
	})
}

func (c *conn) OnBinaryMessage(cb func(message []byte) error) {
	c.OnMessage(func(typ int, message []byte) error {
		if typ == BinaryMessage {
			return cb(message)
		}

		return nil
	})
}

func (c *conn) OnConnect(cb func() error) {
	c.On(event.TypeConnect, eventemitter.HandleFunc(func(payload any) {
		_, ok := payload.(*event.PayloadConnect)
		if !ok {
			c.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := cb(); err != nil {
			c.Emit(event.TypeError, &event.PayloadError{
				Error: err,
			})
		}
	}))
}
