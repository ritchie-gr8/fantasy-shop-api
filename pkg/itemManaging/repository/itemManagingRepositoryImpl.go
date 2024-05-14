package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_itemManagingException "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/exception"
	_itemManagingModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/model"
)

type itemManagingRepositoryImpl struct {
	logger echo.Logger
	db     databases.Database
}

func NewItemManagingRepositoryImpl(logger echo.Logger, db databases.Database) ItemManagingRepository {
	return &itemManagingRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (r *itemManagingRepositoryImpl) CreateItem(item *entities.Item) (*entities.Item, error) {
	result := new(entities.Item)
	if err := r.db.Connect().Create(item).Scan(result).Error; err != nil {
		r.logger.Errorf("failed to create item: %s", err.Error())
		return nil, &_itemManagingException.CreateItem{}
	}

	return result, nil
}

func (r *itemManagingRepositoryImpl) EditItem(itemID uint64, editItemReq *_itemManagingModel.EditItemReq) (uint64, error) {
	if err := r.db.Connect().Model(&entities.Item{}).Where(
		"id = ?", itemID,
	).Updates(editItemReq).Error; err != nil {
		r.logger.Errorf("failed to edit item: %s", err.Error())
		return 0, &_itemManagingException.EditItem{ItemID: itemID}
	}

	return itemID, nil
}

func (r *itemManagingRepositoryImpl) ArchiveItem(itemID uint64) error {
	if err := r.db.Connect().Table(
		"items",
	).Where(
		"id = ?", itemID,
	).Update(
		"is_archive", true,
	).Error; err != nil {
		r.logger.Errorf("failed to archive item: %s", err.Error())
		return &_itemManagingException.ArchiveItem{ItemID: itemID}
	}

	return nil
}
