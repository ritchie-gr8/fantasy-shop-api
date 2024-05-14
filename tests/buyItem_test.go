package tests

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_inventoryRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/repository"
	_itemShopException "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
	_itemShopService "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/service"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	_playerCoinRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBuyItemSuccess(t *testing.T) {
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
	itemShopRepoMock.On("FindByID", uint64(1)).Return(&entities.Item{
		ID:          1,
		Name:        "Sword of Tester",
		Price:       1000,
		Description: "A sword that can be used to test the enemy's defense",
		Picture:     "https://www.google.com/sword-of-tester.jpg",
	}, nil)

	playerCoinRepoMock.On("ShowCoin", "P001").Return(&_playerCoinModel.ShowPlayerCoin{
		PlayerID: "P001",
		Coin:     5000,
	}, nil)

	itemShopRepoMock.On("RecordPurchasHistory", tx, &entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        true,
	}).Return(&entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        true,
	}, nil)

	inventoryRepoMock.On("FillInventory", tx, "P001", uint64(1), int(3)).Return([]*entities.Inventory{
		{
			PlayerID: "P001",
			ItemID:   1,
		},
		{
			PlayerID: "P001",
			ItemID:   1,
		},
		{
			PlayerID: "P001",
			ItemID:   1,
		},
	}, nil)

	playerCoinRepoMock.On("AddCoin", tx, &entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   -3000,
	}).Return(&entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   -3000,
	}, nil)

	type args struct {
		label    string
		in       *_itemShopModel.BuyItemReq
		expected *_playerCoinModel.PlayerCoin
	}

	cases := []args{
		{
			label: "Success: buy item",
			in: &_itemShopModel.BuyItemReq{
				PlayerID: "P001",
				ItemID:   1,
				Quantity: 3,
			},
			expected: &_playerCoinModel.PlayerCoin{
				PlayerID: "P001",
				Amount:   -3000,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			result, err := itemShopService.BuyItem(c.in)
			assert.NoError(t, err)
			assert.EqualValues(t, c.expected, result)
		})
	}
}

func TestBuyItemFail(t *testing.T) {
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
	itemShopRepoMock.On("FindByID", uint64(1)).Return(&entities.Item{
		ID:          1,
		Name:        "Sword of Tester",
		Price:       1000,
		Description: "A sword that can be used to test the enemy's defense",
		Picture:     "https://www.google.com/sword-of-tester.jpg",
	}, nil)

	playerCoinRepoMock.On("ShowCoin", "P001").Return(&_playerCoinModel.ShowPlayerCoin{
		PlayerID: "P001",
		Coin:     2000,
	}, nil)

	type args struct {
		label    string
		in       *_itemShopModel.BuyItemReq
		expected error
	}

	cases := []args{
		{
			label: "Fail: buy item",
			in: &_itemShopModel.BuyItemReq{
				PlayerID: "P001",
				ItemID:   1,
				Quantity: 3,
			},
			expected: &_itemShopException.CoinNotEnough{},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			result, err := itemShopService.BuyItem(c.in)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.EqualValues(t, c.expected, err)
		})
	}
}
