package client

import (
	"github.com/go-zoox/websocket/conn"
)

func (c *client) OnConnect(fn func(conn conn.Conn) error) {
	c.cbs.connects = append(c.cbs.connects, fn)
}
