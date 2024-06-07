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
		events map[string]func(*cs.EventPayload)
	}
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
		eventResponse := &cs.Event{}
		if err := eventResponse.Decode(message); err != nil {
			return err
		}

		if fn, ok := c.cbs.events[eventResponse.ID]; ok {
			delete(c.cbs.events, eventResponse.ID)

			err := safe.Do(func() error {
				fn(&eventResponse.Payload)
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to handle event callback: %v", err)
			}
		}

		return nil
	})

	return c, nil
}
