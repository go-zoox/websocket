package event

import "github.com/go-zoox/websocket/conn"

type PayloadConnect struct {
	Conn conn.Conn
}

const TypeConnect = "connect"
