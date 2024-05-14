package repository

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	"gorm.io/gorm"
)

type PlayerCoinRepository interface {
	AddCoin(tx *gorm.DB, coin *entities.PlayerCoin) (*entities.PlayerCoin, error)
	ShowCoin(playerID string) (*_playerCoinModel.ShowPlayerCoin, error)
}
