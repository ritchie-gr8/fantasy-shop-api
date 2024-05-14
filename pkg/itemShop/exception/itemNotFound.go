package exception

import "fmt"

type ItemNotFound struct {
	ItemID uint64
}

func (e *ItemNotFound) Error() string {
	return fmt.Sprintf("failed to find item by ID: %d", e.ItemID)
}
