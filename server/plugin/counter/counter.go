package counter

import (
	"github.com/go-zoox/websocket/server/plugin"
)

const name = "counter"

type Counter struct {
	plugin.Plugin

	opt *Option
	//
	current int64
	total   int64
}

type Option struct {
}

func New(opts ...func(opt *Option)) *Counter {
	opt := &Option{}
	for _, o := range opts {
		o(opt)
	}

	return &Counter{
		opt: opt,
	}
}

func (hb *Counter) Name() string {
	return name
}

func (hb *Counter) CurrentCount() int64 {
	return hb.current
}

func (hb *Counter) TotalCount() int64 {
	return hb.total
}
