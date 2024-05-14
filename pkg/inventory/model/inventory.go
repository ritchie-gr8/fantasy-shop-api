package model

import (
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
)

type (
	Inventory struct {
		Item     *_itemShopModel.Item `json:"item"`
		Quantity uint                 `json:"quantity"`
	}

	CountItemQuantity struct {
		ItemID   uint64
		Quantity uint
	}
)
