package room

import (
	"fmt"

	"github.com/go-zoox/core-utils/safe"
	"github.com/go-zoox/websocket/conn"
)

var rooms = safe.NewMap[string, Room]()

type Room interface {
	Join(client Client) error
	Leave(client Client) error
	//
	Broadcast(handle func(client Client)) error
	//
	ID() string
}

type Client conn.Conn

type room struct {
	id      string
	clients *safe.Map[string, Client]
}

func New(id string) Room {
	return &room{
		id:      id,
		clients: safe.NewMap[string, Client](),
	}
}

func (r *room) ID() string {
	return r.id
}

func (r *room) Join(client Client) error {
	if r.clients.Has(client.ID()) {
		return fmt.Errorf("[room] failed to join: client(conn: %s) already join the room", client.ID())
	}

	return r.clients.Set(client.ID(), client)
}

func (r *room) Leave(client Client) error {
	if !r.clients.Has(client.ID()) {
		return fmt.Errorf("[room] failed to leave: client(conn: %s) already left the room", client.ID())
	}

	return r.clients.Del(client.ID())
}

func (r *room) Broadcast(handle func(client Client)) error {
	r.clients.ForEach(func(s string, c Client) (stop bool) {
		handle(c)
		return false
	})

	return nil
}
