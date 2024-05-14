package service

import (
	_itemManagingModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/model"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
)

type ItemManagingService interface {
	CreateItem(createItemReq *_itemManagingModel.CreateItemReq) (*_itemShopModel.Item, error)
	EditItem(itemID uint64, editItemReq *_itemManagingModel.EditItemReq) (*_itemShopModel.Item, error)
	ArchiveItem(itemID uint64) error
}
