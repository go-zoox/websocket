package counter

import (
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
)

func (c *Counter) Apply(conn conn.Conn) error {
	conn.OnConnect(func() error {
		c.total.Inc(1)
		c.current.Inc(1)

		logger.Infof("[connections] %d/%d] -> ", c.current.Get(), c.total.Get())
		return nil
	})

	conn.OnClose(func(code int, message string) error {
		c.current.Dec(1)

		logger.Infof("[connections: %d/%d] <- ", c.current.Get(), c.total.Get())
		return nil
	})

	return nil
}
