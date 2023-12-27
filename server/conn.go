package server

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
	"github.com/gorilla/websocket"
)

func (s *server) CreateConn(w http.ResponseWriter, r *http.Request) (conn.Conn, error) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	rawConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return conn.New(r.Context(), rawConn, r), nil
}

func (s *server) ServeConn(connIns conn.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			s.ee.Emit(event.TypeError, &event.PayloadError{
				Error: fmt.Errorf("%v", err),
			})
		}
	}()

	rawConn := connIns.Raw()
	defer func() {
		rawConn.Close()
		// s.ee.Stop()
	}()

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
				return nil
			}

			return err
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
