package client

import (
	"context"
	"time"

	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/websocket/conn"
)

type Client interface {
	OnError(func(conn conn.Conn, err error) error)
	//
	//
	OnConnect(func(conn conn.Conn) error)
	OnClose(func(conn conn.Conn) error)
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
	ee  *eventemitter.EventEmitter
}

type Option struct {
	Context context.Context `json:"context"`
	Addr    string          `json:"addr"`
}

func New(opts ...func(opt *Option)) (Client, error) {
	opt := &Option{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(opt)
	}

	client := &client{
		opt: opt,
		ee:  eventemitter.New(),
	}

	// auto listen pong + send ping
	client.OnPong(func(conn conn.Conn, message []byte) error {
		time.Sleep(15 * time.Second)
		return conn.Ping(message)
	})

	// auto ping first
	client.OnConnect(func(conn conn.Conn) error {
		time.Sleep(5 * time.Second)

		return conn.Ping(nil)
	})

	return client, nil
}
