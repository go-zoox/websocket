package server

import (
	"github.com/go-zoox/websocket/conn"
	c "github.com/go-zoox/websocket/conn"
)

func (s *server) OnMessage(fn func(conn conn.Conn, typ int, message []byte) error) {
	s.cbs.messages = append(s.cbs.messages, fn)
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
