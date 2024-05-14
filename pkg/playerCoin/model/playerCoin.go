package model

import "time"

type (
	PlayerCoin struct {
		ID        uint64    `json:"id"`
		PlayerID  string    `json:"playerID"`
		Amount    int64     `json:"amount"`
		CreatedAt time.Time `json:"createdAt"`
	}

	AddCoinReq struct {
		PlayerID string
		Amount   int64 `json:"amount" validate:"required,gt=0"`
	}

	ShowPlayerCoin struct {
		PlayerID string `json:"playerID"`
		Coin     int64  `json:"coin"`
	}
)
