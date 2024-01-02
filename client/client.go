package client

import (
	"context"
	"net/http"

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
	Headers http.Header     `json:"headers"`
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

	// listen server heartbeat (server ping + client pong)
	client.OnPing(func(conn conn.Conn, message []byte) error {
		return conn.Pong(message)
	})

	return client, nil
}
