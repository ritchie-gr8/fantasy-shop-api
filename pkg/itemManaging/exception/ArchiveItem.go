package exception

import "fmt"

type ArchiveItem struct {
	ItemID uint64
}

func (e *ArchiveItem) Error() string {
	return fmt.Sprintf("failed to archive item id: %d", e.ItemID)
}
