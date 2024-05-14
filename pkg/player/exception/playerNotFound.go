package exception

import "fmt"

type PlayerNotFound struct {
	PlayerID string
}

func (e *PlayerNotFound) Error() string {
	return fmt.Sprintf("failed to find player by id: %s", e.PlayerID)
}
