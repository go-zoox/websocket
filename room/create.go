package room

import "fmt"

func Create(id string) error {
	if ok := rooms.Has(id); ok {
		return fmt.Errorf("[room] %s already exist, cannot create", id)
	}

	rooms.Set(id, New(id))
	return nil
}
