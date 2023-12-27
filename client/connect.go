package client

import (
	"fmt"
	"io"

	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
	"github.com/gorilla/websocket"
)

func (c *client) OnConnect(fn func(conn conn.Conn) error) {
	c.ee.On(event.TypeConnect, eventemitter.HandleFunc(func(payload any) {
		p, ok := payload.(*event.PayloadConnect)
		if !ok {
			c.ee.Emit(event.TypeError, event.ErrInvalidPayload)
			return
		}

		if err := fn(p.Conn); err != nil {
			c.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  p.Conn,
				Error: err,
			})
		}
	}))
}

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
		c.ee.Emit(event.TypePing, &event.PayloadPing{
			Conn:    connIns,
			Message: []byte(appData),
		})
		return nil
	})
	rawConn.SetPongHandler(func(appData string) error {
		c.ee.Emit(event.TypePong, &event.PayloadPong{
			Conn:    connIns,
			Message: []byte(appData),
		})
		return nil
	})

	c.ee.Emit(event.TypeConnect, &event.PayloadConnect{
		Conn: connIns,
	})

	go c.handleConn(connIns, rawConn)

	return nil
}

func (c *client) handleConn(connIns conn.Conn, rawConn *websocket.Conn) {
	for {
		mt, message, err := rawConn.ReadMessage()
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				c.ee.Emit(event.TypeClose, &event.PayloadClose{
					Conn:    connIns,
					Code:    v.Code,
					Message: v.Text,
				})

				return
			}

			c.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  connIns,
				Error: err,
			})
			return
		}

		c.ee.Emit(event.TypeMessage, &event.PayloadMessage{
			Type:    mt,
			Message: message,
		})
	}
}
