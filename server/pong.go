package server

import (
	"github.com/go-zoox/websocket/conn"
)

func (s *server) OnPong(fn func(conn conn.Conn, message []byte) error) {
	s.cbs.pongs = append(s.cbs.pongs, fn)
}
