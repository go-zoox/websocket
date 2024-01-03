package client

import (
	"github.com/go-zoox/websocket/conn"
)

func (c *client) OnClose(cb func(conn conn.Conn, code int, message string) error) {
	c.cbs.closes = append(c.cbs.closes, cb)
}
