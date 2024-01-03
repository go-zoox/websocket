package server

import (
	"fmt"
	"net/http"

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
			Message: []byte(appData),
		})
		return nil
	})
	rawConn.SetPongHandler(func(appData string) error {
		conn.Emit(event.TypePong, &event.PayloadPong{
			Message: []byte(appData),
		})
		return nil
	})

	// event::error
	for _, cb := range s.cbs.errors {
		func(cb func(conn connClass.Conn, err error) error) {
			conn.OnError(func(err error) error {
				return cb(conn, err)
			})
		}(cb)
	}

	// event::close
	for _, cb := range s.cbs.closes {
		func(cb func(conn connClass.Conn, code int, message string) error) {
			conn.OnClose(func(code int, message string) error {
				return cb(conn, code, message)
			})
		}(cb)
	}

	// event::ping
	for _, cb := range s.cbs.pings {
		func(cb func(conn connClass.Conn, message []byte) error) {
			conn.OnPing(func(message []byte) error {
				return cb(conn, message)
			})
		}(cb)
	}

	// event::pong
	for _, cb := range s.cbs.pongs {
		func(cb func(conn connClass.Conn, message []byte) error) {
			conn.OnPong(func(appData []byte) error {
				return cb(conn, appData)
			})
		}(cb)
	}

	// event::message
	for _, cb := range s.cbs.messages {
		func(cb func(conn connClass.Conn, typ int, message []byte) error) {
			conn.OnMessage(func(typ int, message []byte) error {
				return cb(conn, typ, message)
			})
		}(cb)
	}

	// event::connect
	for _, cb := range s.cbs.connects {
		func(cb func(conn connClass.Conn) error) {
			conn.OnConnect(func() error {
				return cb(conn)
			})
		}(cb)
	}

	return conn, nil
}

func (s *server) ServeConn(conn connClass.Conn) {
	defer conn.Close()

	defer func() {
		if err := recover(); err != nil {
			s.ee.Emit(event.TypeError, &event.PayloadError{
				Error: fmt.Errorf("%v", err),
			})
		}
	}()

	conn.Emit(event.TypeConnect, &event.PayloadConnect{})

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				conn.Emit(event.TypeClose, &event.PayloadClose{
					Code:    v.Code,
					Message: v.Text,
				})
				return
			}

			s.ee.Emit(event.TypeError, &event.PayloadError{
				Error: err,
			})
			return
		}

		// do not hold the message reader
		go conn.Emit(event.TypeMessage, &event.PayloadMessage{
			Type:    mt,
			Message: message,
		})
	}
}
