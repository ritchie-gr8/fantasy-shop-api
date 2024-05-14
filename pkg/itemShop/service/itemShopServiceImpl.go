package service

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_inventoryRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/repository"
	_itemShopException "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	_playerCoinRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/repository"
)

type itemShopServiceImpl struct {
	itemShopRepository _itemShopRepository.ItemShopRepository
	playerCoinRepo     _playerCoinRepo.PlayerCoinRepository
	inventoryRepo      _inventoryRepo.InventoryRepository
	logger             echo.Logger
}

func NewItemShopServiceImpl(
	itemShopRepository _itemShopRepository.ItemShopRepository,
	playerCoinRepo _playerCoinRepo.PlayerCoinRepository,
	inventoryRepo _inventoryRepo.InventoryRepository,
	logger echo.Logger,
) ItemShopService {
	return &itemShopServiceImpl{
		itemShopRepository: itemShopRepository,
		playerCoinRepo:     playerCoinRepo,
		inventoryRepo:      inventoryRepo,
		logger:             logger,
	}
}

func (s *itemShopServiceImpl) GetItems(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error) {
	items, err := s.itemShopRepository.GetItems(itemFilter)
	if err != nil {
		return nil, err
	}

	totalItemsNum, err := s.itemShopRepository.CountItems(itemFilter)
	if err != nil {
		return nil, err
	}

	totalPage := s.totalPageCalculation(totalItemsNum, itemFilter.Size)

	return s.toItemResultResponse(items, itemFilter.Page, totalPage), nil
}

func (s *itemShopServiceImpl) BuyItem(buyReq *_itemShopModel.BuyItemReq) (*_playerCoinModel.PlayerCoin, error) {
	// get item detail
	itemEntity, err := s.itemShopRepository.FindByID(buyReq.ItemID)
	if err != nil {
		return nil, err
	}

	// calculate total price
	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), buyReq.Quantity)
	if err := s.checkPlayerCoin(buyReq.PlayerID, totalPrice); err != nil {
		return nil, err
	}

	// add purchase history
	tx := s.itemShopRepository.TransactionBegin()
	purchaseHistory, err := s.itemShopRepository.RecordPurchasHistory(
		tx, &entities.PurchaseHistory{
			PlayerID:        buyReq.PlayerID,
			ItemID:          buyReq.ItemID,
			ItemName:        itemEntity.Name,
			ItemDescription: itemEntity.Description,
			ItemPrice:       itemEntity.Price,
			ItemPicture:     itemEntity.Picture,
			Quantity:        buyReq.Quantity,
			IsBuying:        true,
		},
	)
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("recorded purchase history ID: %d", purchaseHistory.ID)

	// deduct player coin
	playerCoin, err := s.playerCoinRepo.AddCoin(
		tx, &entities.PlayerCoin{
			PlayerID: buyReq.PlayerID,
			Amount:   -totalPrice,
		},
	)
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Info("deducted player coin: %d", playerCoin.Amount)

	// add item into player's inventory
	inventoryEntity, err := s.inventoryRepo.FillInventory(
		tx, buyReq.PlayerID, buyReq.ItemID, int(buyReq.Quantity),
	)
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("fill playerID%s inventory: %d", buyReq.PlayerID, len(inventoryEntity))

	//commit transaction
	if err := s.itemShopRepository.TransactionCommit(tx); err != nil {
		return nil, err
	}

	return playerCoin.ToPlayerCoinModel(), nil
}

func (s *itemShopServiceImpl) SellItem(sellReq *_itemShopModel.SellItemReq) (*_playerCoinModel.PlayerCoin, error) {
	// get item detail
	itemEntity, err := s.itemShopRepository.FindByID(sellReq.ItemID)
	if err != nil {
		return nil, err
	}

	// calculate total price
	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), sellReq.Quantity)
	totalPrice = totalPrice / 2
	if err := s.checkPlayerItem(
		sellReq.PlayerID,
		itemEntity.ID,
		sellReq.Quantity,
	); err != nil {
		return nil, err
	}

	// add purchase history
	tx := s.itemShopRepository.TransactionBegin()
	purchaseHistory, err := s.itemShopRepository.RecordPurchasHistory(
		tx, &entities.PurchaseHistory{
			PlayerID:        sellReq.PlayerID,
			ItemID:          sellReq.ItemID,
			ItemName:        itemEntity.Name,
			ItemDescription: itemEntity.Description,
			ItemPrice:       itemEntity.Price,
			ItemPicture:     itemEntity.Picture,
			Quantity:        sellReq.Quantity,
			IsBuying:        false,
		},
	)
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("recorded purchase history ID: %d", purchaseHistory.ID)

	// add player coin
	playerCoin, err := s.playerCoinRepo.AddCoin(
		tx, &entities.PlayerCoin{
			PlayerID: sellReq.PlayerID,
			Amount:   totalPrice,
		},
	)
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Info("deducted player coin: %d", playerCoin.Amount)

	// add item into player's inventory
	err = s.inventoryRepo.RemoveItemFromInventory(
		tx, sellReq.PlayerID, sellReq.ItemID, int(sellReq.Quantity),
	)
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("remove playerID:%s itemID:%d from inventory: %d qty",
		sellReq.PlayerID, sellReq.ItemID, sellReq.Quantity)

	//commit transaction
	if err := s.itemShopRepository.TransactionCommit(tx); err != nil {
		return nil, err
	}

	return playerCoin.ToPlayerCoinModel(), nil
}

func (s *itemShopServiceImpl) totalPageCalculation(totalItem, size int64) int64 {
	totalPage := totalItem / size
	if totalItem%size != 0 {
		totalPage++
	}

	return totalPage
}

func (s *itemShopServiceImpl) toItemResultResponse(
	items []*entities.Item, page, totalPage int64) *_itemShopModel.ItemResult {

	modeledItems := make([]*_itemShopModel.Item, 0)
	for _, item := range items {
		modeledItems = append(modeledItems, item.ToItemModel())
	}

	return &_itemShopModel.ItemResult{
		Items: modeledItems,
		Paginate: _itemShopModel.PaginateResult{
			Page:      page,
			TotalPage: totalPage,
		},
	}
}

func (s *itemShopServiceImpl) totalPriceCalculation(item *_itemShopModel.Item, qty uint) int64 {
	return int64(item.Price) * int64(qty)
}

func (s *itemShopServiceImpl) checkPlayerCoin(playerID string, totalPrice int64) error {
	playerCoin, err := s.playerCoinRepo.ShowCoin(playerID)
	if err != nil {
		return err
	}

	if playerCoin.Coin < totalPrice {
		s.logger.Error("failed because player coin is not enough")
		return &_itemShopException.CoinNotEnough{}
	}

	return nil
}

func (s *itemShopServiceImpl) checkPlayerItem(playerID string, itemID uint64, sellQty uint) error {
	itemQty := s.inventoryRepo.CountPlayerItem(playerID, itemID)
	if itemQty < int64(sellQty) {
		s.logger.Error("failed because player does not have enough item")
		return &_itemShopException.ItemNotEnough{ItemID: itemID}
	}

	return nil
}
