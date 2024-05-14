package exception

type AddCoin struct{}

func (e *AddCoin) Error() string {
	return "failed to add coin"
}
