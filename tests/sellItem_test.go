package tests

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_inventoryRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/repository"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"

	_itemShopException "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/exception"
	_itemShopRepository "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
	_itemShopService "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/service"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	_playerCoinRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSellItemSuccess(t *testing.T) {
	itemShopRepoMock := new(_itemShopRepository.ItemShopRepositoryMock)
	playerCoinRepoMock := new(_playerCoinRepo.PlayerCoinRepositoryMock)
	inventoryRepoMock := new(_inventoryRepo.InventoryRepositoryMock)
	logger := echo.New().Logger

	itemShopService := _itemShopService.NewItemShopServiceImpl(
		itemShopRepoMock,
		playerCoinRepoMock,
		inventoryRepoMock,
		logger,
	)

	tx := &gorm.DB{}
	itemShopRepoMock.On("TransactionBegin").Return(tx)
	itemShopRepoMock.On("TransactionCommit", tx).Return(nil)
	itemShopRepoMock.On("TransactionRollback", tx).Return(nil)

	inventoryRepoMock.On("CountPlayerItem", "P001", uint64(1)).Return(int64(3), nil)

	itemShopRepoMock.On("FindByID", uint64(1)).Return(&entities.Item{
		ID:          1,
		Name:        "Sword of Tester",
		Price:       1000,
		Description: "A sword that can be used to test the enemy's defense",
		Picture:     "https://www.google.com/sword-of-tester.jpg",
	}, nil)

	itemShopRepoMock.On("RecordPurchasHistory", tx, &entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        false,
	}).Return(&entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        false,
	}, nil)

	playerCoinRepoMock.On("AddCoin", tx, &entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   1500,
	}).Return(&entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   1500,
	}, nil)

	inventoryRepoMock.On("RemoveItemFromInventory", tx, "P001", uint64(1), 3).Return(nil)

	type args struct {
		label    string
		in       *_itemShopModel.SellItemReq
		expected *_playerCoinModel.PlayerCoin
	}

	cases := []args{
		{
			"Success: sell item",
			&_itemShopModel.SellItemReq{
				PlayerID: "P001",
				ItemID:   1,
				Quantity: 3,
			},
			&_playerCoinModel.PlayerCoin{
				PlayerID: "P001",
				Amount:   1500,
			},
		},
	}

	for _, c := range cases {
		result, err := itemShopService.SellItem(c.in)
		assert.NoError(t, err)
		assert.EqualValues(t, c.expected, result)
	}
}

func TestSellItemFail(t *testing.T) {
	itemShopRepoMock := new(_itemShopRepository.ItemShopRepositoryMock)
	playerCoinRepoMock := new(_playerCoinRepo.PlayerCoinRepositoryMock)
	inventoryRepoMock := new(_inventoryRepo.InventoryRepositoryMock)
	logger := echo.New().Logger

	itemShopService := _itemShopService.NewItemShopServiceImpl(
		itemShopRepoMock,
		playerCoinRepoMock,
		inventoryRepoMock,
		logger,
	)

	tx := &gorm.DB{}
	itemShopRepoMock.On("TransactionBegin").Return(tx)
	itemShopRepoMock.On("TransactionCommit", tx).Return(nil)
	itemShopRepoMock.On("TransactionRollback", tx).Return(nil)

	inventoryRepoMock.On("CountPlayerItem", "P001", uint64(1)).Return(int64(2), nil)

	itemShopRepoMock.On("FindByID", uint64(1)).Return(&entities.Item{
		ID:          1,
		Name:        "Sword of Tester",
		Price:       1000,
		Description: "A sword that can be used to test the enemy's defense",
		Picture:     "https://www.google.com/sword-of-tester.jpg",
	}, nil)

	type args struct {
		label    string
		in       *_itemShopModel.SellItemReq
		expected error
	}

	cases := []args{
		{
			"Fail: sell item",
			&_itemShopModel.SellItemReq{
				PlayerID: "P001",
				ItemID:   1,
				Quantity: 3,
			},
			&_itemShopException.ItemNotEnough{ItemID: 1},
		},
	}

	for _, c := range cases {
		result, err := itemShopService.SellItem(c.in)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualValues(t, c.expected, err)
	}
}
