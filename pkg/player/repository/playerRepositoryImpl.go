package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_playerException "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/exception"
)

type playerRepositoryImpl struct {
	logger echo.Logger
	db     databases.Database
}

func NewPlayerRepositoryImpl(logger echo.Logger, db databases.Database) PlayerRepository {
	return &playerRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (r *playerRepositoryImpl) Create(player *entities.Player) (*entities.Player, error) {
	result := new(entities.Player)
	if err := r.db.Connect().Create(player).Scan(result).Error; err != nil {
		r.logger.Errorf("failed to create player: %s", err.Error())
		return nil, &_playerException.CreatePlayer{PlayerID: player.ID}
	}

	return result, nil
}

func (r *playerRepositoryImpl) FindByID(playerID string) (*entities.Player, error) {
	player := new(entities.Player)
	if err := r.db.Connect().Where("id = ?", playerID).First(player).Error; err != nil {
		r.logger.Errorf("failed to find player by ID: %s", err.Error())
		return nil, &_playerException.PlayerNotFound{PlayerID: playerID}
	}

	return player, nil
}
