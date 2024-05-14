package exception

import "fmt"

type CreateAdmin struct {
	AdminID string
}

func (e *CreateAdmin) Error() string {
	return fmt.Sprintf("failed to create admin id: %s", e.AdminID)
}
