package service

import (
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
)

type ItemShopService interface {
	GetItems(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error)
	BuyItem(buyReq *_itemShopModel.BuyItemReq) (*_playerCoinModel.PlayerCoin, error)
	SellItem(sellReq *_itemShopModel.SellItemReq) (*_playerCoinModel.PlayerCoin, error)
}
