package counter

import (
	"github.com/go-zoox/websocket/server/plugin"
)

const name = "counter"

type counter struct {
	opt *Option
	//
	Count int
}

type Option struct {
}

func New(opts ...func(opt *Option)) plugin.Plugin {
	opt := &Option{}
	for _, o := range opts {
		o(opt)
	}

	return &counter{
		opt: opt,
	}
}

func (hb *counter) Name() string {
	return name
}
