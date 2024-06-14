package client

import (
	"context"
	"net/http"
	"time"

	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
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

	SendMessage(message []byte) error
	SendBinaryMessage(message []byte) error
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

	return c, nil
}
