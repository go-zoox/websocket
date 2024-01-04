package counter

import (
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
)

func (c *counter) Apply(conn conn.Conn) error {
	conn.OnConnect(func() error {
		c.Count += 1

		logger.Infof("[connections] %d] -> ", c.Count)
		return nil
	})

	conn.OnClose(func(code int, message string) error {
		c.Count -= 1
		logger.Infof("[connections: %d] <- ", c.Count)
		return nil
	})

	return nil
}
