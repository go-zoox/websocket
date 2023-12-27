package server

import (
	"github.com/go-zoox/websocket/conn"
)

func (s *server) OnPing(fn func(conn conn.Conn, message []byte) error) {
	s.cbs.pings = append(s.cbs.pings, fn)
}
