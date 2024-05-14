package exception

type GetItems struct {
}

func (e *GetItems) Error() string {
	return "failed to get items"
}
