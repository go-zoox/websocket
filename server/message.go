package server

import (
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/websocket/conn"
	c "github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
)

func (s *server) OnMessage(fn func(conn conn.Conn, typ int, message []byte) error) {
	s.ee.On(event.TypeMessage, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadMessage)
		if !ok {
			s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := fn(p.Conn, p.Type, p.Message); err != nil {
			s.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  p.Conn,
				Error: err,
			})
		}
	}))
}

func (s *server) OnTextMessage(cb func(conn conn.Conn, message []byte) error) {
	s.OnMessage(func(conn c.Conn, typ int, message []byte) error {
		if typ == c.TextMessage {
			return cb(conn, message)
		}

		return nil
	})
}

func (s *server) OnBinaryMessage(cb func(conn conn.Conn, message []byte) error) {
	s.OnMessage(func(conn c.Conn, typ int, message []byte) error {
		if typ == c.BinaryMessage {
			return cb(conn, message)
		}

		return nil
	})
}
