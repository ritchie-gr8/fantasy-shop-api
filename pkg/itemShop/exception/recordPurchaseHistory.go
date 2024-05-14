package exception

type RecordPurchasHistory struct{}

func (e *RecordPurchasHistory) Error() string {
	return "failed to record purchas history"
}
