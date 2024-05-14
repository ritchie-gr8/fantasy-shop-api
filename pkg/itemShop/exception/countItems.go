package exception

type CountItems struct{}

func (e *CountItems) Error() string {
	return "failed to count items"
}
