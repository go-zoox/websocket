package server

import (
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event/cs"
)

// Event registers a event with handler
func (s *server) Event(name string, fn func(conn conn.Conn, payload *cs.EventPayload, callback func(*cs.EventPayload)) error) {
	s.cbs.events[name] = append(s.cbs.events[name], fn)
}
