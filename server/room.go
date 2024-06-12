package server

import "github.com/go-zoox/websocket/room"

func (s *server) Room(id string) room.Room {
	return room.Get(id)
}
