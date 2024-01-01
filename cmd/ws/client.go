package main

import (
	"fmt"
	"time"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/websocket/client"
	"github.com/go-zoox/websocket/conn"
)

func Client() *cli.Command {
	return &cli.Command{
		Name:  "client",
		Usage: "websocket client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dsn",
				Usage: "websocket server address",
				Value: "ws://localhost:9000",
			},
			&cli.StringFlag{
				Name:  "message",
				Usage: "websocket server message",
				Value: "hello world",
			},
		},
		Action: func(c *cli.Context) error {
			wc, err := client.New(func(opt *client.Option) {
				opt.Addr = c.String("dsn")
			})
			if err != nil {
				return err
			}

			wc.OnError(func(conn conn.Conn, err error) error {
				fmt.Println("websocket client error:", err)
				return nil
			})

			wc.OnClose(func(conn conn.Conn) error {
				fmt.Println("websocket client closed")
				return nil
			})

			wc.OnPing(func(conn conn.Conn, message []byte) error {
				fmt.Println("websocket client on ping:", string(message))
				return nil
			})

			wc.OnPong(func(conn conn.Conn, message []byte) error {
				fmt.Println("websocket client on pong:", string(message))
				return nil
			})

			wc.OnConnect(func(conn conn.Conn) error {
				fmt.Println("websocket client connected")
				// conn.Ping(nil)

				i := 0
				for {
					if i > 100 {
						continue
					}
					i++

					time.Sleep(1 * time.Second)
					err = conn.WriteTextMessage([]byte(fmt.Sprintf("%s - %d", c.String("message"), time.Now().Unix())))
					if err != nil {
						return err
					}
				}
				// return nil
			})

			// if err := wc.Connect(); err != nil {
			// 	return err
			// }

			// if err := wc.SendTextMessage([]byte("hello world")); err != nil {
			// 	return err
			// }

			// i := 0
			// for {
			// 	if i > 100 {
			// 		continue
			// 	}
			// 	i++

			// 	time.Sleep(1 * time.Second)
			// 	err = wc.SendTextMessage([]byte(fmt.Sprintf("hello world %d", time.Now().Unix())))
			// 	if err != nil {
			// 		return err
			// 	}
			// }

			// wc.Ping()

			// time.Sleep(1e6 * time.Second)

			// wc.Close()

			if err := wc.Connect(); err != nil {
				return err
			}

			for {
				select {}
			}

			return nil
		},
	}
}
