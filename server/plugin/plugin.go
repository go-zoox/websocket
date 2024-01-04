package plugin

import "github.com/go-zoox/websocket/conn"

type Plugin interface {
	Name() string
	Apply(conn conn.Conn) error
}
