package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-zoox/logger"
	"github.com/go-zoox/safe"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event/cs"
)

type Client interface {
	OnError(func(conn conn.Conn, err error) error)
	//
	//
	OnConnect(func(conn conn.Conn) error)
	OnClose(func(conn conn.Conn, code int, message string) error)
	//
	OnMessage(func(conn conn.Conn, typ int, message []byte) error)
	//
	OnPing(func(conn conn.Conn, message []byte) error)
	OnPong(func(conn conn.Conn, message []byte) error)
	//
	OnTextMessage(func(conn conn.Conn, message []byte) error)
	OnBinaryMessage(func(conn conn.Conn, message []byte) error)

	//
	Connect() error
	Close() error
	Reconnect() error

	//
	Event(typ string, payload cs.EventPayload, callback func(err error, payload cs.EventPayload), opts ...EventOption) error
}

type client struct {
	conn conn.Conn

	opt *Option

	cbs struct {
		errors   []func(conn conn.Conn, err error) error
		connects []func(conn conn.Conn) error
		closes   []func(conn conn.Conn, code int, message string) error
		messages []func(conn conn.Conn, typ int, message []byte) error
		pings    []func(conn conn.Conn, message []byte) error
		pongs    []func(conn conn.Conn, message []byte) error
		//
		events map[string]EventCallback
	}
}

type EventCallback struct {
	Callback    func(err error, payload cs.EventPayload)
	IsSubscribe bool
}

type Option struct {
	Context context.Context `json:"context"`
	Addr    string          `json:"addr"`
	Headers http.Header     `json:"headers"`

	//
	ConnectTimeout time.Duration `json:"connect_timeout"`
}

func New(opts ...func(opt *Option)) (Client, error) {
	opt := &Option{
		Context:        context.Background(),
		ConnectTimeout: 30 * time.Second,
	}
	for _, o := range opts {
		o(opt)
	}

	c := &client{
		opt: opt,
	}
	c.cbs.events = make(map[string]EventCallback)

	// listen server heartbeat (server ping + client pong)
	c.OnPing(func(conn conn.Conn, message []byte) error {
		logger.Debugf("[heartbeat][interval][ping] receive heartbeat <-")

		if err := conn.Pong(message); err != nil {
			logger.Errorf("[heartbeat][interval] fail to send heartbeat: %v", err)
			return err
		}

		logger.Debugf("[heartbeat][interval][pong] send heartbeat ->")
		return nil
	})

	// event
	c.OnTextMessage(func(conn conn.Conn, message []byte) error {
		go func() {
			event := &cs.Event{}
			if err := event.Decode(message); err != nil {
				logger.Errorf("[event] failed to decode: %s (message: %s)", err, string(message))
				// return err
			}

			if fn, ok := c.cbs.events[event.ID]; ok {
				if !fn.IsSubscribe {
					delete(c.cbs.events, event.ID)
				}

				err := safe.Do(func() error {
					var err error
					if event.Error != "" {
						err = fmt.Errorf(event.Error)
					}

					fn.Callback(err, event.Payload)
					return nil
				})
				if err != nil {
					logger.Errorf("failed to handle event callback: %v", err)
				}
			}
		}()

		return nil
	})

	return c, nil
}
