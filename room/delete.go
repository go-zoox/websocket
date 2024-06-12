package room

import "fmt"

func Delete(id string) error {
	if ok := rooms.Has(id); !ok {
		return fmt.Errorf("[room] %s already deleted, cannot delete", id)
	}

	return rooms.Del(id)
}
