package counter

import (
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
)

func (c *Counter) Apply(conn conn.Conn) error {
	conn.OnConnect(func() error {
		c.current += 1
		c.total += 1

		logger.Infof("[connections] %d/%d] -> ", c.current, c.total)
		return nil
	})

	conn.OnClose(func(code int, message string) error {
		c.current -= 1
		logger.Infof("[connections: %d/%d] <- ", c.current, c.total)
		return nil
	})

	return nil
}
