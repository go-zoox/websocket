package server

import (
	"github.com/go-zoox/websocket/conn"
)

func (s *server) OnClose(fn func(conn conn.Conn, code int, message string) error) {
	s.cbs.closes = append(s.cbs.closes, fn)
}
