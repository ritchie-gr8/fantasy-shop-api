package repository

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
	"gorm.io/gorm"
)

type ItemShopRepository interface {
	TransactionBegin() *gorm.DB
	TransactionRollback(tx *gorm.DB) error
	TransactionCommit(tx *gorm.DB) error
	GetItems(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error)
	CountItems(itemFilter *_itemShopModel.ItemFilter) (int64, error)
	FindByID(itemID uint64) (*entities.Item, error)
	FindByIDList(itemIDs []uint64) ([]*entities.Item, error)
	RecordPurchasHistory(tx *gorm.DB, purchaseEntity *entities.PurchaseHistory) (*entities.PurchaseHistory, error)
}
