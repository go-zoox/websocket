package client

import "time"

func (c *client) Reconnect() error {
	time.Sleep(1 * time.Second)
	return c.Connect()
}
