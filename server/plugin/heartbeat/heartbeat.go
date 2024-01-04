package heartbeat

import (
	"time"

	"github.com/go-zoox/websocket/server/plugin"
)

const name = "heartbeat"

type heartbeat struct {
	opt *Option
}

type Option struct {
	Interval time.Duration
	Timeout  time.Duration
}

func New(opts ...func(opt *Option)) plugin.Plugin {
	opt := &Option{}
	for _, o := range opts {
		o(opt)
	}

	return &heartbeat{
		opt: opt,
	}
}
