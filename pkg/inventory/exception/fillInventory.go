package exception

import "fmt"

type FillInventory struct {
	PlayerID string
	ItemID   uint64
}

func (e *FillInventory) Error() string {
	return fmt.Sprintf(
		"failed to fill inventory playerID: %s, itemID: %d",
		e.PlayerID, e.ItemID,
	)
}
