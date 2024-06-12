package room

func Join(client Client, id string) error {
	return rooms.Get(id).Join(client)
}
