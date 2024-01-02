package server

import (
	"net/http"
	"time"

	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
	connClass "github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
)

type Server interface {
	Run(addr string) error
	//
	OnError(func(conn conn.Conn, err error) error)
	//
	OnConnect(func(conn conn.Conn) error)
	OnClose(func(conn conn.Conn) error)
	//
	OnMessage(func(conn conn.Conn, typ int, message []byte) error)
	//
	OnPing(func(conn conn.Conn, message []byte) error)
	OnPong(func(conn conn.Conn, message []byte) error)
	//
	OnTextMessage(func(conn conn.Conn, message []byte) error)
	OnBinaryMessage(func(conn conn.Conn, message []byte) error)
	//
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	//
	CreateConn(w http.ResponseWriter, r *http.Request) (conn.Conn, error)
	ServeConn(connIns conn.Conn)
}

type Option struct {
}

type server struct {
	opt *Option
	//
	ee *eventemitter.EventEmitter
	//
	cbs struct {
		errors   []func(conn conn.Conn, err error) error
		connects []func(conn conn.Conn) error
		closes   []func(conn conn.Conn) error
		messages []func(conn conn.Conn, typ int, message []byte) error
		pings    []func(conn conn.Conn, message []byte) error
		pongs    []func(conn conn.Conn, message []byte) error
	}
}

func New(opts ...func(opt *Option)) (Server, error) {
	opt := &Option{}
	for _, o := range opts {
		o(opt)
	}

	s := &server{
		opt: opt,
		ee:  eventemitter.New(),
	}

	// event::error
	for _, cb := range s.cbs.errors {
		func(cb func(conn connClass.Conn, err error) error) {
			s.ee.On(event.TypeError, eventemitter.HandleFunc(func(payload any) {
				p, ok := payload.(*event.PayloadError)
				if !ok {
					// s.ee.Emit(event.TypeError, event.ErrInvalidPayload)
					logger.Errorf("invalid payload: %v", payload)
					return
				}

				if err := cb(p.Conn, p.Error); err != nil {
					logger.Errorf("failed to handle error: %v", err)
				}
			}))
		}(cb)
	}

	// @TODO auto listen ping + sennd pong
	s.OnPong(func(conn conn.Conn, message []byte) error {
		conn.Get("heartbeat").(chan struct{}) <- struct{}{}
		return nil
	})

	//
	s.OnConnect(func(conn conn.Conn) error {
		ch := make(chan struct{})
		conn.Set("heartbeat", ch)

		// heartbeat listener
		go func() {
			logger.Debugf("heartbeat started")

			timer := time.NewTicker(15 * time.Second)
			for {
				select {
				case <-conn.Context().Done():
					logger.Debugf("[heartbeat][listener] context done => cancel")
					return
				case <-timer.C:
					logger.Errorf("[heartbeat][listener] fail to listen heartbeat")
					close(ch)
					go conn.Close()
					return
				case <-ch:
					logger.Debugf("[heartbeat][listener] receive heartbeat <-")
					timer.Reset(15 * time.Second)
				}
			}
		}()

		// heartbeat sender
		go func() {
			for {
				select {
				case <-conn.Context().Done():
					logger.Debugf("[heartbeat][sender] context done => cancel")
					return
				case <-time.After(10 * time.Second):
					logger.Debugf("[heartbeat][sender] send heartbeat ->")
					if err := conn.Ping(nil); err != nil {
						logger.Errorf("[heartbeat][sender] fail to send heartbeat: %v", err)
						return
					}
				}
			}
		}()

		return nil
	})

	return s, nil
}
