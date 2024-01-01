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

	// @TODO auto listen ping + sennd pong
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

	// event::error
	for _, cb := range s.cbs.errors {
		func(cb func(conn connClass.Conn, err error) error) {
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
		}(cb)
	}

	// event::connect
	for _, cb := range s.cbs.connects {
		func(cb func(conn connClass.Conn) error) {
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
		}(cb)
	}

	// event::close
	for _, cb := range s.cbs.closes {
		func(cb func(conn connClass.Conn) error) {
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
		}(cb)
	}

	// event::ping
	for _, cb := range s.cbs.pings {
		func(cb func(conn connClass.Conn, message []byte) error) {
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
		}(cb)

	}

	// event::pong
	for _, cb := range s.cbs.pongs {
		func(cb func(conn connClass.Conn, message []byte) error) {
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
		}(cb)
	}

	// event::message
	for _, cb := range s.cbs.messages {
		func(cb func(conn connClass.Conn, typ int, message []byte) error) {
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
		}(cb)
	}

	return conn, nil
}

func (s *server) ServeConn(conn connClass.Conn) {
	defer conn.Close()

	defer func() {
		if err := recover(); err != nil {
			conn.Emit(event.TypeError, &event.PayloadError{
				Conn:  conn,
				Error: fmt.Errorf("%v", err),
			})
		}
	}()

	conn.Emit(event.TypeConnect, &event.PayloadConnect{
		Conn: conn,
	})

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				conn.Emit(event.TypeClose, &event.PayloadClose{
					Conn:    conn,
					Code:    v.Code,
					Message: v.Text,
				})
				return
			}

			conn.Emit(event.TypeError, &event.PayloadError{
				Conn:  conn,
				Error: err,
			})
			return
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
