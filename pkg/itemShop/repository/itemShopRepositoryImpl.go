package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_itemShopException "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
	"gorm.io/gorm"
)

type itemShopRepositoryImpl struct {
	logger echo.Logger
	db     databases.Database
}

func NewItemShopRepositoryImpl(logger echo.Logger, db databases.Database) ItemShopRepository {
	return &itemShopRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (r *itemShopRepositoryImpl) TransactionBegin() *gorm.DB {
	tx := r.db.Connect()
	return tx.Begin()
}

func (r *itemShopRepositoryImpl) TransactionRollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *itemShopRepositoryImpl) TransactionCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *itemShopRepositoryImpl) GetItems(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error) {
	items := make([]*entities.Item, 0)

	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false) // select * from items
	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%")
	}

	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%")
	}

	offset := int((itemFilter.Page - 1) * itemFilter.Size)
	limit := int(itemFilter.Size)

	if err := query.Offset(offset).Limit(limit).Find(&items).Order("id asc").Error; err != nil {
		r.logger.Errorf("failed to get items: %s", err.Error())
		return nil, &_itemShopException.GetItems{}
	}

	return items, nil
}

func (r *itemShopRepositoryImpl) CountItems(itemFilter *_itemShopModel.ItemFilter) (int64, error) {

	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false) // select * from items
	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%")
	}

	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%")
	}

	var count int64

	if err := query.Count(&count).Error; err != nil {
		r.logger.Errorf("failed to count items: %s", err.Error())
		return -1, &_itemShopException.CountItems{}
	}

	return count, nil
}

func (r *itemShopRepositoryImpl) FindByID(itemID uint64) (*entities.Item, error) {
	item := new(entities.Item)

	if err := r.db.Connect().First(item, itemID).Error; err != nil {
		r.logger.Errorf("failed to find item by ID: %s", err.Error())
		return nil, &_itemShopException.ItemNotFound{}
	}

	return item, nil
}

func (r *itemShopRepositoryImpl) FindByIDList(itemIDs []uint64) ([]*entities.Item, error) {
	items := make([]*entities.Item, 0)
	if err := r.db.Connect().Model(
		&entities.Item{},
	).Where(
		"id in ?", itemIDs,
	).Find(&items).Error; err != nil {
		r.logger.Errorf("failed to find item by ID list: %s", err.Error())
		return nil, &_itemShopException.GetItems{}
	}

	return items, nil
}

func (r *itemShopRepositoryImpl) RecordPurchasHistory(
	tx *gorm.DB,
	purchaseEntity *entities.PurchaseHistory,
) (*entities.PurchaseHistory, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	purchaseHistory := new(entities.PurchaseHistory)
	if err := conn.Create(purchaseEntity).Error; err != nil {
		r.logger.Errorf("failed to record purchase history: %s", err.Error())
		return nil, &_itemShopException.RecordPurchasHistory{}
	}

	return purchaseHistory, nil
}
