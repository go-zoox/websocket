package client

import (
	"fmt"
	"io"

	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
	"github.com/gorilla/websocket"
)

func (c *client) Connect() error {
	rawConn, response, err := websocket.DefaultDialer.DialContext(c.opt.Context, c.opt.Addr, c.opt.Headers)
	if err != nil {
		if response == nil || response.Body == nil {
			return fmt.Errorf("failed to connect at %s (error: %s)", c.opt.Addr, err)
		}

		body, errB := io.ReadAll(response.Body)
		if errB != nil {
			return fmt.Errorf("failed to connect at %s (status: %s, error: %s)", c.opt.Addr, response.Status, err)
		}

		return fmt.Errorf("failed to connect at %s (status: %d, response: %s, error: %v)", c.opt.Addr, response.StatusCode, string(body), err)
	}
	// defer conn.Close()

	connIns := conn.New(c.opt.Context, rawConn, nil)
	c.conn = connIns

	rawConn.SetPingHandler(func(appData string) error {
		c.conn.Emit(event.TypePing, &event.PayloadPing{
			Message: []byte(appData),
		})
		return nil
	})
	rawConn.SetPongHandler(func(appData string) error {
		c.conn.Emit(event.TypePong, &event.PayloadPong{
			Message: []byte(appData),
		})
		return nil
	})

	// event::error
	for _, cb := range c.cbs.errors {
		func(cb func(conn conn.Conn, err error) error) {
			c.conn.OnError(func(err error) error {
				return cb(c.conn, err)
			})
		}(cb)
	}

	// event::close
	for _, cb := range c.cbs.closes {
		func(cb func(conn conn.Conn, code int, message string) error) {
			c.conn.OnClose(func(code int, message string) error {
				return cb(c.conn, code, message)
			})
		}(cb)
	}

	// event::ping
	for _, cb := range c.cbs.pings {
		func(cb func(conn conn.Conn, message []byte) error) {
			c.conn.OnPing(func(message []byte) error {
				return cb(c.conn, message)
			})
		}(cb)
	}

	// event::pong
	for _, cb := range c.cbs.pongs {
		func(cb func(conn conn.Conn, message []byte) error) {
			c.conn.OnPong(func(appData []byte) error {
				return cb(c.conn, appData)
			})
		}(cb)
	}

	// event::message
	for _, cb := range c.cbs.messages {
		func(cb func(conn conn.Conn, typ int, message []byte) error) {
			c.conn.OnMessage(func(typ int, message []byte) error {
				return cb(c.conn, typ, message)
			})
		}(cb)
	}

	// event::connect
	for _, cb := range c.cbs.connects {
		func(cb func(conn conn.Conn) error) {
			c.conn.OnConnect(func() error {
				return cb(c.conn)
			})
		}(cb)
	}

	c.conn.Emit(event.TypeConnect, &event.PayloadConnect{})

	go c.handleConn(connIns, rawConn)

	return nil
}

func (c *client) handleConn(connIns conn.Conn, rawConn *websocket.Conn) {
	for {
		mt, message, err := rawConn.ReadMessage()
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				c.conn.Emit(event.TypeClose, &event.PayloadClose{
					Code:    v.Code,
					Message: v.Text,
				})

				return
			}

			c.conn.Emit(event.TypeError, &event.PayloadError{
				Error: err,
			})
			return
		}

		c.conn.Emit(event.TypeMessage, &event.PayloadMessage{
			Type:    mt,
			Message: message,
		})
	}
}

func (c *client) Close() error {
	return c.conn.Close()
}
