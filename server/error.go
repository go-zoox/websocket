package server

import (
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
)

func (s *server) OnError(fn func(conn conn.Conn, err error) error) {
	s.ee.On(event.TypeError, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadError)
		if !ok {
			s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := fn(p.Conn, p.Error); err != nil {
			logger.Errorf("failed to handle error: %v", err)
		}
	}))
}
