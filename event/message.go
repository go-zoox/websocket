package event

import "github.com/go-zoox/websocket/conn"

type PayloadMessage struct {
	Conn    conn.Conn
	Type    int
	Message []byte
}

const TypeMessage = "[internal] message"
