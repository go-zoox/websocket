package server

import (
	"net/http"
	"time"

	"github.com/go-zoox/websocket/conn"
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
	ServeConn(connIns conn.Conn) error
}

type Option struct {
}

type server struct {
	opt *Option
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
	}

	// @TODO auto listen ping + sennd pong
	s.OnPing(func(conn conn.Conn, message []byte) error {
		time.Sleep(time.Second * 3)
		return conn.Pong(message)
	})

	return s, nil
}
