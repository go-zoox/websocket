package server

import (
	"github.com/go-zoox/websocket/conn"
)

func (s *server) OnConnect(fn func(conn conn.Conn) error) {
	s.cbs.connects = append(s.cbs.connects, fn)
}
