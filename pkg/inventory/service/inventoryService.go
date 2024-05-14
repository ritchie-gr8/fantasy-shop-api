package service

import (
	_inventoryModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/model"
)

type InventoryService interface {
	GetInventory(playerID string) ([]*_inventoryModel.Inventory, error)
}
