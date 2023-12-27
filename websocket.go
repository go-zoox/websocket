package websocket

import (
	"github.com/go-zoox/websocket/client"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/server"
)

type Conn = conn.Conn

type Server = server.Server

type Client = client.Client

func NewServer(opts ...func(opt *server.Option)) (Server, error) {
	return server.New(opts...)
}

func NewClient(opts ...func(opt *client.Option)) (Client, error) {
	return client.New(opts...)
}
