package exception

type CreateItem struct{}

func (e *CreateItem) Error() string {
	return "failed to create item"
}
