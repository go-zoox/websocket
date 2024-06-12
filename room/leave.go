package room

func Leave(client Client, id string) error {
	return rooms.Get(id).Leave(client)
}
