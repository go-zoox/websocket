package event

import "github.com/go-zoox/websocket/conn"

type PayloadPong struct {
	Conn    conn.Conn
	Message []byte
}

const TypePong = "[internal] pong"
