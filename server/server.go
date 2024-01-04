package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
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
	ServeConn(connIns conn.Conn)
}

type Option struct {
	HearbeatInterval time.Duration
	HeartbeatTimeout time.Duration
}

type server struct {
	opt *Option
	//
	ee *eventemitter.EventEmitter
	//
	cbs struct {
		errors   []func(conn conn.Conn, err error) error
		connects []func(conn conn.Conn) error
		closes   []func(conn conn.Conn, code int, message string) error
		messages []func(conn conn.Conn, typ int, message []byte) error
		pings    []func(conn conn.Conn, message []byte) error
		pongs    []func(conn conn.Conn, message []byte) error
	}
}

func New(opts ...func(opt *Option)) (Server, error) {
	opt := &Option{
		HearbeatInterval: 25 * time.Second,
		HeartbeatTimeout: 15 * time.Second,
	}
	for _, o := range opts {
		o(opt)
	}

	if opt.HearbeatInterval < opt.HeartbeatTimeout {
		return nil, fmt.Errorf("heartbeat interval must be greater than heartbeat timeout")
	}

	s := &server{
		opt: opt,
		ee:  eventemitter.New(),
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

		// heartbeat
		go func() {
			time.After(opt.HearbeatInterval)
			for {
				select {
				case <-conn.Context().Done():
					logger.Debugf("[heartbeat][interval] context done => cancel")
					return
				case <-time.After(opt.HearbeatInterval):
					logger.Debugf("[heartbeat][interval] send heartbeat ->")
					if err := conn.Ping(nil); err != nil {
						logger.Errorf("[heartbeat][interval] fail to send heartbeat: %v", err)
						close(ch)
						go conn.Close()
						return
					}

					select {
					case <-conn.Context().Done():
						logger.Debugf("[heartbeat][timeout] context done => cancel")
						return
					case <-time.After(opt.HeartbeatTimeout):
						logger.Errorf("[heartbeat][timeout] fail to listen heartbeat")
						close(ch)
						go conn.Close()
						return
					case <-ch:
						logger.Debugf("[heartbeat][timeout] receive heartbeat <-")
					}
				}
			}
		}()

		return nil
	})

	return s, nil
}
