package model

type (
	BuyItemReq struct {
		PlayerID string
		ItemID   uint64 `json:"itemID" validate:"required,gt=0"`
		Quantity uint   `json:"quantity" validate:"required,gt=0"`
	}

	SellItemReq struct {
		PlayerID string
		ItemID   uint64 `json:"itemID" validate:"required,gt=0"`
		Quantity uint   `json:"quantity" validate:"required,gt=0"`
	}
)
