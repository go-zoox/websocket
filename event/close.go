package event

import "github.com/go-zoox/websocket/conn"

type PayloadClose struct {
	Conn conn.Conn
	//
	Code    int
	Message string
}

const TypeClose = "close"
