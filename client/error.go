package client

import (
	"github.com/go-zoox/websocket/conn"
)

func (c *client) OnError(fn func(conn conn.Conn, err error) error) {
	c.cbs.errors = append(c.cbs.errors, fn)
}
