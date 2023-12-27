package websocket

import (
	"github.com/go-zoox/websocket/client"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/server"
)

type Conn = conn.Conn

type Server = server.Server

type ServerOption = server.Option

type Client = client.Client

type ClientOption = client.Option

func NewServer(opts ...func(opt *ServerOption)) (Server, error) {
	return server.New(opts...)
}

func NewClient(opts ...func(opt *ClientOption)) (Client, error) {
	return client.New(opts...)
}
