package counter

import (
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
)

func (c *Counter) Apply(conn conn.Conn) error {
	conn.OnConnect(func() error {
		c.total.Set(c.total.Get() + 1)
		c.current.Set(c.current.Get() + 1)

		logger.Infof("[connections] %d/%d] -> ", c.current.Get(), c.total.Get())
		return nil
	})

	conn.OnClose(func(code int, message string) error {
		c.current.Set(c.current.Get() - 1)

		logger.Infof("[connections: %d/%d] <- ", c.current.Get(), c.total.Get())
		return nil
	})

	return nil
}
