package server

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/websocket/event"
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connIns, err := s.CreateConn(w, r)
	if err != nil {
		connIns.Emit(event.TypeError, &event.PayloadError{
			Error: fmt.Errorf("%v", err),
		})
		return
	}

	s.ServeConn(connIns)
}
