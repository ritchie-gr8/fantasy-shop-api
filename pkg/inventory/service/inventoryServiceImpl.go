package service

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_inventoryModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/model"
	_inventoryRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/repository"
	_itemShopRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
)

type inventoryServiceImpl struct {
	inventoryRepo _inventoryRepo.InventoryRepository
	itemShopRepo  _itemShopRepo.ItemShopRepository
}

func NewInventoryServiceImpl(
	inventoryRepo _inventoryRepo.InventoryRepository,
	itemShopRepo _itemShopRepo.ItemShopRepository,
) InventoryService {
	return &inventoryServiceImpl{
		inventoryRepo: inventoryRepo,
		itemShopRepo:  itemShopRepo,
	}
}

func (s *inventoryServiceImpl) GetInventory(playerID string) ([]*_inventoryModel.Inventory, error) {
	inventory, err := s.inventoryRepo.GetInventory(playerID)
	if err != nil {
		return nil, err
	}

	itemListWithQuantity := s.getUniqueItemWithQuantityList(inventory)

	return s.mapGetInventoryResponse(itemListWithQuantity), nil
}

func (s *inventoryServiceImpl) getUniqueItemWithQuantityList(
	inventoryEntities []*entities.Inventory,
) []_inventoryModel.CountItemQuantity {
	list := make([]_inventoryModel.CountItemQuantity, 0)
	itemMap := make(map[uint64]uint)

	for _, item := range inventoryEntities {
		itemMap[item.ItemID]++
	}

	for itemID, quantity := range itemMap {
		list = append(list, _inventoryModel.CountItemQuantity{
			ItemID:   itemID,
			Quantity: quantity,
		})
	}

	return list
}

func (s *inventoryServiceImpl) mapGetInventoryResponse(
	itemListWithQuantity []_inventoryModel.CountItemQuantity,
) []*_inventoryModel.Inventory {

	itemIDList := s.getItemByID(itemListWithQuantity)
	itemList, err := s.itemShopRepo.FindByIDList(itemIDList)
	if err != nil {
		return []*_inventoryModel.Inventory{}
	}

	result := make([]*_inventoryModel.Inventory, 0)
	itemMapWithQuantity := s.getItemMapWithQuantity(itemListWithQuantity)

	for _, item := range itemList {
		result = append(result, &_inventoryModel.Inventory{
			Item:     item.ToItemModel(),
			Quantity: itemMapWithQuantity[item.ID],
		})
	}

	return result
}

func (s *inventoryServiceImpl) getItemByID(
	itemListWithQuantity []_inventoryModel.CountItemQuantity,
) []uint64 {

	itemIDList := make([]uint64, 0)

	for _, item := range itemListWithQuantity {
		itemIDList = append(itemIDList, item.ItemID)
	}

	return itemIDList
}

func (s *inventoryServiceImpl) getItemMapWithQuantity(
	uniqueItemList []_inventoryModel.CountItemQuantity,
) map[uint64]uint {
	result := make(map[uint64]uint)

	for _, item := range uniqueItemList {
		result[item.ItemID] = item.Quantity
	}
	return result
}
