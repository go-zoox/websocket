package client

import (
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
)

func (s *client) OnPong(fn func(conn conn.Conn, message []byte) error) {
	s.ee.On(event.TypePong, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadPong)
		if !ok {
			s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := fn(p.Conn, p.Message); err != nil {
			s.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  p.Conn,
				Error: err,
			})
		}
	}))
}
