package server

import (
	"github.com/go-zoox/websocket/conn"
)

func (s *server) Use(fn func(conn conn.Conn, next func())) {
	s.cbs.middlewares = append(s.cbs.middlewares, fn)
}
