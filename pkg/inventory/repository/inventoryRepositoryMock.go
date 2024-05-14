package repository

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type InventoryRepositoryMock struct {
	mock.Mock
}

func (m *InventoryRepositoryMock) FillInventory(
	tx *gorm.DB, playerID string,
	itemID uint64, limit int) ([]*entities.Inventory, error) {
	args := m.Called(tx, playerID, itemID, limit)
	return args.Get(0).([]*entities.Inventory), args.Error(1)
}

func (m *InventoryRepositoryMock) RemoveItemFromInventory(
	tx *gorm.DB, playerID string,
	itemID uint64, limit int) error {
	args := m.Called(tx, playerID, itemID, limit)
	return args.Error(0)
}

func (m *InventoryRepositoryMock) CountPlayerItem(
	playerID string, itemID uint64) int64 {
	args := m.Called(playerID, itemID)
	return args.Get(0).(int64)
}

func (m *InventoryRepositoryMock) GetInventory(
	playerID string) ([]*entities.Inventory, error) {
	args := m.Called(playerID)
	return args.Get(0).([]*entities.Inventory), args.Error(1)
}
