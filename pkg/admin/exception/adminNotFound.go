package exception

import "fmt"

type AdminNotFound struct {
	AdminID string
}

func (e *AdminNotFound) Error() string {
	return fmt.Sprintf("failed to find admin by id: %s", e.AdminID)
}
