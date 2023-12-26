package server

import (
	"net"
	"net/http"

	"github.com/go-zoox/logger"
)

func (s *server) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	server := &http.Server{
		Addr:    addr,
		Handler: s,
	}

	logger.Infof("websocket server is running on %s", addr)

	return server.Serve(listener)
}
