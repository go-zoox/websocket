package server

import (
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
)

func (s *server) OnClose(fn func(conn conn.Conn) error) {
	s.ee.On(event.TypeClose, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadClose)
		if !ok {
			s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := fn(p.Conn); err != nil {
			s.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  p.Conn,
				Error: err,
			})
		}
	}))
}
