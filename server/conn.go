package server

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/logger"
	connClass "github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
	"github.com/gorilla/websocket"
)

func (s *server) CreateConn(w http.ResponseWriter, r *http.Request) (connClass.Conn, error) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	rawConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	conn := connClass.New(r.Context(), rawConn, r)

	// event::error
	for _, cb := range s.cbs.errors {
		conn.On(event.TypeError, eventemitter.HandleFunc(func(payload any) {
			p, ok := payload.(*event.PayloadError)
			if !ok {
				conn.Emit(event.TypeError, event.ErrInvalidPayload)
				return
			}

			if err := cb(p.Conn, p.Error); err != nil {
				logger.Errorf("failed to handle error: %v", err)
			}
		}))
	}

	// event::connect
	for _, cb := range s.cbs.connects {
		conn.On(event.TypeConnect, eventemitter.HandleFunc(func(payload any) {
			p, ok := payload.(*event.PayloadConnect)
			if !ok {
				conn.Emit(event.TypeError, event.ErrInvalidPayload)
				return
			}

			if err := cb(p.Conn); err != nil {
				conn.Emit(event.TypeError, &event.PayloadError{
					Conn:  p.Conn,
					Error: err,
				})
			}
		}))
	}

	// event::close
	for _, cb := range s.cbs.closes {
		conn.On(event.TypeClose, eventemitter.HandleFunc(func(payload any) {
			p, ok := payload.(*event.PayloadClose)
			if !ok {
				conn.Emit(event.TypeError, event.ErrInvalidPayload)
				return
			}

			if err := cb(p.Conn); err != nil {
				conn.Emit(event.TypeError, &event.PayloadError{
					Conn:  p.Conn,
					Error: err,
				})
			}
		}))
	}

	// event::ping
	for _, cb := range s.cbs.pings {
		conn.On(event.TypePing, eventemitter.HandleFunc(func(payload any) {
			p, ok := payload.(*event.PayloadPing)
			if !ok {
				conn.Emit(event.TypeError, event.ErrInvalidPayload)
				return
			}

			if err := cb(p.Conn, p.Message); err != nil {
				conn.Emit(event.TypeError, &event.PayloadError{
					Conn:  p.Conn,
					Error: err,
				})
			}
		}))
	}

	// event::pong
	for _, cb := range s.cbs.pongs {
		conn.On(event.TypePong, eventemitter.HandleFunc(func(payload any) {
			p, ok := payload.(*event.PayloadPong)
			if !ok {
				conn.Emit(event.TypeError, event.ErrInvalidPayload)
				return
			}

			if err := cb(p.Conn, p.Message); err != nil {
				conn.Emit(event.TypeError, &event.PayloadError{
					Conn:  p.Conn,
					Error: err,
				})
			}
		}))
	}

	// event::message
	for _, cb := range s.cbs.messages {
		conn.On(event.TypeMessage, eventemitter.HandleFunc(func(payload any) {
			p, ok := payload.(*event.PayloadMessage)
			if !ok {
				conn.Emit(event.TypeError, event.ErrInvalidPayload)
				return
			}

			if err := cb(p.Conn, p.Type, p.Message); err != nil {
				conn.Emit(event.TypeError, &event.PayloadError{
					Conn:  p.Conn,
					Error: err,
				})
			}
		}))
	}

	return conn, nil
}

func (s *server) ServeConn(conn connClass.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			conn.Emit(event.TypeError, &event.PayloadError{
				Error: fmt.Errorf("%v", err),
			})
		}
	}()

	rawConn := conn.Raw()
	defer func() {
		rawConn.Close()
		conn.Close()
	}()

	conn.Emit(event.TypeConnect, &event.PayloadConnect{
		Conn: conn,
	})

	rawConn.SetPingHandler(func(appData string) error {
		conn.Emit(event.TypePing, &event.PayloadPing{
			Conn:    conn,
			Message: []byte(appData),
		})
		return nil
	})
	rawConn.SetPongHandler(func(appData string) error {
		conn.Emit(event.TypePong, &event.PayloadPong{
			Conn:    conn,
			Message: []byte(appData),
		})
		return nil
	})

	for {
		mt, message, err := rawConn.ReadMessage()
		fmt.Println("mt", mt, "message", message, "err", err)
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				conn.Emit(event.TypeClose, &event.PayloadClose{
					Conn:    conn,
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
						conn.Emit(event.TypeError, &event.PayloadError{
							Conn:  conn,
							Error: v,
						})
					} else {
						conn.Emit(event.TypeError, &event.PayloadError{
							Conn:  conn,
							Error: fmt.Errorf("%v", err),
						})
					}
				}
			}()

			conn.Emit(event.TypeMessage, &event.PayloadMessage{
				Conn:    conn,
				Type:    mt,
				Message: message,
			})
		}()
	}
}
