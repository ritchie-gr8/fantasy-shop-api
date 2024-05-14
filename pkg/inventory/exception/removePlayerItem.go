package exception

import "fmt"

type RemovePlayerItem struct {
	ItemID uint64
}

func (e *RemovePlayerItem) Error() string {
	return fmt.Sprintf("failed to remove item itemID: %d", e.ItemID)
}
