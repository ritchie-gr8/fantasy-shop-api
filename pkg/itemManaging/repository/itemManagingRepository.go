package repository

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_itemManagingModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/model"
)

type ItemManagingRepository interface {
	CreateItem(item *entities.Item) (*entities.Item, error)
	EditItem(itemID uint64, editItemReq *_itemManagingModel.EditItemReq) (uint64, error)
	ArchiveItem(itemID uint64) error
}
