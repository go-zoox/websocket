package server

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
	"github.com/gorilla/websocket"
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			s.ee.Emit(event.TypeError, &event.PayloadError{
				Error: fmt.Errorf("%v", err),
			})
		}
	}()

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	rawConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.ee.Emit(event.TypeError, &event.PayloadError{
			Error: fmt.Errorf("%v", err),
		})
		return
	}

	connIns := conn.New(r.Context(), rawConn)
	defer rawConn.Close()

	s.ee.Emit(event.TypeConnect, &event.PayloadConnect{
		Conn: connIns,
	})

	rawConn.SetPingHandler(func(appData string) error {
		s.ee.Emit(event.TypePing, &event.PayloadPing{
			Conn:    connIns,
			Message: []byte(appData),
		})
		return nil
	})
	rawConn.SetPongHandler(func(appData string) error {
		s.ee.Emit(event.TypePong, &event.PayloadPong{
			Conn:    connIns,
			Message: []byte(appData),
		})
		return nil
	})

	for {
		mt, message, err := rawConn.ReadMessage()
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				s.ee.Emit(event.TypeClose, &event.PayloadClose{
					Conn:    connIns,
					Code:    v.Code,
					Message: v.Text,
				})
				return
			}

			s.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  connIns,
				Error: err,
			})
			return
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					if v, ok := err.(error); ok {
						s.ee.Emit(event.TypeError, &event.PayloadError{
							Conn:  connIns,
							Error: v,
						})
					} else {
						s.ee.Emit(event.TypeError, &event.PayloadError{
							Conn:  connIns,
							Error: fmt.Errorf("%v", err),
						})
					}
				}
			}()

			s.ee.Emit(event.TypeMessage, &event.PayloadMessage{
				Conn:    connIns,
				Type:    mt,
				Message: message,
			})
		}()
	}
}
