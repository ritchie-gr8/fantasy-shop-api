package exception

type CoinNotEnough struct{}

func (e *CoinNotEnough) Error() string {
	return "failed because coin is not enough"
}
