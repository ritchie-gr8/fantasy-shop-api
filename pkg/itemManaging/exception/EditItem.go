package exception

import "fmt"

type EditItem struct {
	ItemID uint64
}

func (e *EditItem) Error() string {
	return fmt.Sprintf("failed to edit item id: %d", e.ItemID)
}
