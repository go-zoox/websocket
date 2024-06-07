package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event"
	"github.com/go-zoox/websocket/server/plugin"
	"github.com/go-zoox/websocket/server/plugin/counter"
	"github.com/go-zoox/websocket/server/plugin/heartbeat"
)

type Server interface {
	Run(addr string) error
	//
	OnError(func(conn conn.Conn, err error) error)
	//
	OnConnect(func(conn conn.Conn) error)
	OnClose(func(conn conn.Conn, code int, message string) error)
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
	ServeConn(conn conn.Conn)
	//
	Plugin(plugin plugin.Plugin) error
	//
	Use(middleware func(conn conn.Conn, next func()))
}

type Option struct {
	HeartbeatInterval time.Duration
	HeartbeatTimeout  time.Duration
}

type server struct {
	opt *Option
	//
	ee eventemitter.EventEmitter
	//
	cbs struct {
		errors   []func(conn conn.Conn, err error) error
		connects []func(conn conn.Conn) error
		closes   []func(conn conn.Conn, code int, message string) error
		messages []func(conn conn.Conn, typ int, message []byte) error
		pings    []func(conn conn.Conn, message []byte) error
		pongs    []func(conn conn.Conn, message []byte) error
		//
		middlewares []func(conn conn.Conn, next func())
		//
		events map[string][]func(conn conn.Conn, payload *EventPayload) error
	}
	//
	plugins map[string]plugin.Plugin
}

func New(opts ...func(opt *Option)) (Server, error) {
	opt := &Option{
		HeartbeatInterval: 25 * time.Second,
		HeartbeatTimeout:  15 * time.Second,
	}
	for _, o := range opts {
		o(opt)
	}

	if opt.HeartbeatInterval < opt.HeartbeatTimeout {
		return nil, fmt.Errorf("heartbeat interval must be greater than heartbeat timeout")
	}

	s := &server{
		opt:     opt,
		ee:      eventemitter.New(),
		plugins: make(map[string]plugin.Plugin),
	}

	s.ee.On(event.TypeError, eventemitter.HandleFunc(func(payload any) {
		if err, ok := payload.(*event.PayloadError); ok {
			logger.Errorf("[server] internal error: %s", err.Error)
		}
	}))

	s.Plugin(heartbeat.New(func(o *heartbeat.Option) {
		o.Interval = opt.HeartbeatInterval
		o.Timeout = opt.HeartbeatTimeout
	}))

	s.Plugin(counter.New())

	s.OnTextMessage(func(c conn.Conn, message []byte) error {
		eventX := &Event{}
		if err := eventX.Decode(message); err != nil {
			return err
		}

		if fns, ok := s.cbs.events[eventX.Type]; ok {
			for _, fn := range fns {
				go (func(fn func(c conn.Conn, payload *EventPayload) error) {
					if err := fn(c, &eventX.Payload); err != nil {
						s.ee.Emit(event.TypeError, &event.PayloadError{
							Error: fmt.Errorf("failed to handle event(type: %s): %v", eventX.Type, err),
						})
					}
				})(fn)
			}
		} else if !ok {
			s.ee.Emit(event.TypeError, &event.PayloadError{
				Error: fmt.Errorf("supported event type: %s", eventX.Type),
			})
		}

		return nil
	})

	return s, nil
}
