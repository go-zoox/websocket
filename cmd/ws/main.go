package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/websocket"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "ws",
		Usage:   "websocket server/client",
		Version: websocket.Version,
	})

	app.Register("client", Client())
	app.Register("server", Server())

	app.Run()
}
