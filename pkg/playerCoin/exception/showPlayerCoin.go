package exception

type ShowPlayerCoin struct{}

func (e *ShowPlayerCoin) Error() string {
	return "failed to show player coin"
}
