package exception

import "fmt"

type CreatePlayer struct {
	PlayerID string
}

func (e *CreatePlayer) Error() string {
	return fmt.Sprintf("failed to create player id: %s", e.PlayerID)
}
