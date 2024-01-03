package client

import (
	"github.com/go-zoox/websocket/conn"
)

func (c *client) OnPong(fn func(conn conn.Conn, message []byte) error) {
	c.cbs.pongs = append(c.cbs.pongs, fn)
}
