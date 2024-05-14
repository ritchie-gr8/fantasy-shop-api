package repository

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	FillInventory(tx *gorm.DB, playerID string, itemID uint64, limit int) ([]*entities.Inventory, error)
	RemoveItemFromInventory(tx *gorm.DB, playerID string, itemID uint64, limit int) error
	CountPlayerItem(playerID string, itemID uint64) int64
	GetInventory(playerID string) ([]*entities.Inventory, error)
}
