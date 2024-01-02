package main

import (
	"fmt"
	"time"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/server"
)

func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "websocket server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dsn",
				Usage: "websocket server address",
			},
		},
		Action: func(c *cli.Context) error {
			s, err := server.New()
			if err != nil {
				return err
			}

			s.OnError(func(conn conn.Conn, err error) error {
				fmt.Println("websocket server error:", err)
				return nil
			})

			s.OnConnect(func(conn conn.Conn) error {
				fmt.Println("websocket server connected")
				return nil
			})

			s.OnClose(func(conn conn.Conn) error {
				fmt.Println("websocket server closed")
				return nil
			})

			s.OnPing(func(conn conn.Conn, message []byte) error {
				fmt.Println("websocket server on ping:", string(message))
				// conn.Pong(nil)
				// conn.Ping(nil)
				return nil
			})

			// s.OnPong(func(conn conn.Conn, message []byte) error {
			// 	fmt.Println("websocket server on pong:", string(message))
			// 	// conn.Ping(message)
			// 	return nil
			// })

			s.OnMessage(func(conn conn.Conn, typ int, message []byte) error {
				fmt.Println("websocket server message:", typ, string(message))
				time.Sleep(2 * time.Second)
				return nil
			})

			return s.Run(":9000")
		},
	}
}
