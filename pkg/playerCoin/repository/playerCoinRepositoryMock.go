package repository

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type PlayerCoinRepositoryMock struct {
	mock.Mock
}

func (m *PlayerCoinRepositoryMock) AddCoin(tx *gorm.DB, coin *entities.PlayerCoin) (*entities.PlayerCoin, error) {
	args := m.Called(tx, coin)
	return args.Get(0).(*entities.PlayerCoin), args.Error(1)
}

func (m *PlayerCoinRepositoryMock) ShowCoin(playerID string) (*_playerCoinModel.ShowPlayerCoin, error) {
	args := m.Called(playerID)
	return args.Get(0).(*_playerCoinModel.ShowPlayerCoin), args.Error(1)
}
