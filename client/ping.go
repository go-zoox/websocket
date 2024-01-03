package client

import (
	"github.com/go-zoox/websocket/conn"
)

func (c *client) OnPing(fn func(conn conn.Conn, message []byte) error) {
	c.cbs.pings = append(c.cbs.pings, fn)
}
