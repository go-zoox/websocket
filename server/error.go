package server

import (
	"github.com/go-zoox/websocket/conn"
)

func (s *server) OnError(fn func(conn conn.Conn, err error) error) {
	s.cbs.errors = append(s.cbs.errors, fn)
}
