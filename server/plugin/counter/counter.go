package counter

import (
	"github.com/go-zoox/core-utils/safe"
	"github.com/go-zoox/websocket/server/plugin"
)

const name = "counter"

type Counter struct {
	plugin.Plugin

	opt *Option
	//
	current safe.Int64
	total   safe.Int64
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
	return hb.current.Get()
}

func (hb *Counter) TotalCount() int64 {
	return hb.total.Get()
}
