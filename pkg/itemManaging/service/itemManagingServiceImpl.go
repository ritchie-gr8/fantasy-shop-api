package service

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_itemManagingModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/model"
	_itemManagingRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/repository"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
	_itemShopRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
)

type itemManagingServiceImpl struct {
	itemManagingRepository _itemManagingRepo.ItemManagingRepository
	itemShopRepository     _itemShopRepo.ItemShopRepository
}

func NewItemShopServiceImpl(
	itemManagingRepository _itemManagingRepo.ItemManagingRepository,
	itemShopRepository _itemShopRepo.ItemShopRepository,
) ItemManagingService {
	return &itemManagingServiceImpl{
		itemManagingRepository: itemManagingRepository,
		itemShopRepository:     itemShopRepository,
	}
}

func (s *itemManagingServiceImpl) CreateItem(
	createItemReq *_itemManagingModel.CreateItemReq,
) (*_itemShopModel.Item, error) {

	item := &entities.Item{
		Name:        createItemReq.Name,
		Description: createItemReq.Description,
		Picture:     createItemReq.Picture,
		Price:       createItemReq.Price,
	}

	result, err := s.itemManagingRepository.CreateItem(item)
	if err != nil {
		return nil, err
	}

	return result.ToItemModel(), nil
}

func (s *itemManagingServiceImpl) EditItem(
	itemID uint64,
	editItemReq *_itemManagingModel.EditItemReq,
) (*_itemShopModel.Item, error) {
	itemID, err := s.itemManagingRepository.EditItem(itemID, editItemReq)
	if err != nil {
		return nil, err
	}

	result, err := s.itemShopRepository.FindByID(itemID)
	if err != nil {
		return nil, err
	}

	return result.ToItemModel(), nil
}

func (s *itemManagingServiceImpl) ArchiveItem(itemID uint64) error {
	return s.itemManagingRepository.ArchiveItem(itemID)
}
