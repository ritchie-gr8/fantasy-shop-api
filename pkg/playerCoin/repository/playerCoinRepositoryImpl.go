package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_playerCoinException "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/exception"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	"gorm.io/gorm"
)

type playerCoinRepositoryImpl struct {
	logger echo.Logger
	db     databases.Database
}

func NewPlayerCoinRepositoryImpl(
	logger echo.Logger,
	db databases.Database,
) PlayerCoinRepository {
	return &playerCoinRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (r *playerCoinRepositoryImpl) AddCoin(tx *gorm.DB, coin *entities.PlayerCoin) (*entities.PlayerCoin, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	playerCoin := new(entities.PlayerCoin)

	if err := conn.Create(coin).Scan(playerCoin).Error; err != nil {
		r.logger.Errorf("failed to add player coin: %s", err.Error())
		return nil, &_playerCoinException.AddCoin{}
	}

	return playerCoin, nil
}

func (r *playerCoinRepositoryImpl) ShowCoin(playerID string) (*_playerCoinModel.ShowPlayerCoin, error) {
	showCoinModel := new(_playerCoinModel.ShowPlayerCoin)

	if err := r.db.Connect().Model(
		&entities.PlayerCoin{},
	).Where(
		"player_id = ?", playerID,
	).Select(
		"player_id, sum(amount) as coin",
	).Group(
		"player_id",
	).Scan(showCoinModel).Error; err != nil {
		r.logger.Errorf("failed to show player coin: %s", err.Error())
		return nil, &_playerCoinException.ShowPlayerCoin{}
	}

	return showCoinModel, nil
}
