package event

import "github.com/go-zoox/websocket/conn"

type PayloadPing struct {
	Conn    conn.Conn
	Message []byte
}

const TypePing = "[internal] ping"
