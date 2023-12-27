package client

import (
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
	rawConn, _, err := websocket.DefaultDialer.DialContext(c.opt.Context, c.opt.Addr, nil)
	if err != nil {
		return err
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

	for {
		mt, message, err := rawConn.ReadMessage()
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				c.ee.Emit(event.TypeClose, &event.PayloadClose{
					Conn:    connIns,
					Code:    v.Code,
					Message: v.Text,
				})

				return err
			}

			c.ee.Emit(event.TypeError, &event.PayloadError{
				Conn:  connIns,
				Error: err,
			})
			return err
		}

		c.ee.Emit(event.TypeMessage, &event.PayloadMessage{
			Type:    mt,
			Message: message,
		})
	}
}
