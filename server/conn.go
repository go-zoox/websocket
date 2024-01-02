package server

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/eventemitter"
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

	// event::connect
	for _, cb := range s.cbs.connects {
		func(cb func(conn connClass.Conn) error) {
			conn.On(event.TypeConnect, eventemitter.HandleFunc(func(payload any) {
				p, ok := payload.(*event.PayloadConnect)
				if !ok {
					s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
					return
				}

				if err := cb(p.Conn); err != nil {
					s.ee.Emit(event.TypeError, &event.PayloadError{
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
					s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
					return
				}

				if err := cb(p.Conn); err != nil {
					s.ee.Emit(event.TypeError, &event.PayloadError{
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
					s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
					return
				}

				if err := cb(p.Conn, p.Message); err != nil {
					s.ee.Emit(event.TypeError, &event.PayloadError{
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
					s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
					return
				}

				if err := cb(p.Conn, p.Message); err != nil {
					s.ee.Emit(event.TypeError, &event.PayloadError{
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
					s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
					return
				}

				if err := cb(p.Conn, p.Type, p.Message); err != nil {
					s.ee.Emit(event.TypeError, &event.PayloadError{
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
			s.ee.Emit(event.TypeError, &event.PayloadError{
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

			s.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  conn,
				Error: err,
			})
			return
		}

		conn.Emit(event.TypeMessage, &event.PayloadMessage{
			Conn:    conn,
			Type:    mt,
			Message: message,
		})
	}
}
