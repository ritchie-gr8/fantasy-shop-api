package exception

import "fmt"

type FindPlayerItem struct {
	PlayerID string
}

func (e *FindPlayerItem) Error() string {
	return fmt.Sprintf("failed to find player item playerID: %s", e.PlayerID)
}
