package room

func Get(id string) Room {
	if ok := rooms.Has(id); !ok {
		rooms.Set(id, New(id))
	}

	return rooms.Get(id)
}
