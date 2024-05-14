package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_inventoryException "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/exception"
	"gorm.io/gorm"
)

type inventoryRepositoryImpl struct {
	logger echo.Logger
	db     databases.Database
}

func NewInventoryRepositoryImpl(
	logger echo.Logger,
	db databases.Database,
) InventoryRepository {
	return &inventoryRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (r *inventoryRepositoryImpl) FillInventory(
	tx *gorm.DB, playerID string, itemID uint64, limit int,
) ([]*entities.Inventory, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	result := make([]*entities.Inventory, 0)

	for range limit {
		result = append(result, &entities.Inventory{
			PlayerID: playerID,
			ItemID:   itemID,
		})
	}

	if err := conn.Create(result).Error; err != nil {
		r.logger.Errorf("failed to fill inventory: %s", err.Error())
		return nil, &_inventoryException.FillInventory{
			PlayerID: playerID,
			ItemID:   itemID,
		}
	}

	return result, nil
}

func (r *inventoryRepositoryImpl) RemoveItemFromInventory(
	tx *gorm.DB,
	playerID string, itemID uint64, limit int) error {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	inventory, err := r.findItemInPlayerInventoryByID(playerID, itemID, limit)
	if err != nil {
		return err
	}

	for _, item := range inventory {
		item.IsArchive = true
		if err := conn.Model(
			&entities.Inventory{},
		).Where(
			"id = ?", item.ID,
		).Updates(
			item,
		).Error; err != nil {
			tx.Rollback()
			r.logger.Errorf("failed to remove item from player inventory itemID: %s", err.Error())
			return &_inventoryException.RemovePlayerItem{
				ItemID: item.ID,
			}
		}
	}

	return nil
}

func (r *inventoryRepositoryImpl) CountPlayerItem(playerID string, itemID uint64) int64 {
	var count int64

	if err := r.db.Connect().Model(
		&entities.Inventory{},
	).Where(
		"player_id = ? and item_id = ? and is_archive = ?",
		playerID, itemID, false,
	).Count(&count).Error; err != nil {
		r.logger.Errorf("failed to count player item: %s", err.Error())
		return -1
	}

	return count
}

func (r *inventoryRepositoryImpl) GetInventory(playerID string) ([]*entities.Inventory, error) {
	result := make([]*entities.Inventory, 0)

	if err := r.db.Connect().Where(
		"player_id = ? and is_archive = ?",
		playerID, false,
	).Find(&result).Error; err != nil {
		r.logger.Errorf("failed to get player inventory: %s", err.Error())
		return nil, &_inventoryException.FindPlayerItem{
			PlayerID: playerID,
		}
	}

	return result, nil
}

func (r *inventoryRepositoryImpl) findItemInPlayerInventoryByID(
	playerID string,
	itemID uint64,
	limit int,
) ([]*entities.Inventory, error) {
	result := make([]*entities.Inventory, 0)
	if err := r.db.Connect().Where(
		"player_id = ? and item_id = ? and is_archive = ?",
		playerID, itemID, false,
	).Limit(
		limit,
	).Find(&result).Error; err != nil {
		r.logger.Errorf("failed to find item in player inventory by ID: %s", err.Error())
		return nil, &_inventoryException.RemovePlayerItem{
			ItemID: itemID,
		}
	}

	return result, nil
}
