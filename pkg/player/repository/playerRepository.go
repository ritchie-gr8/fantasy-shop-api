package repository

import "github.com/ritchie-gr8/fantasy-shop-api/entities"

type PlayerRepository interface {
	Create(player *entities.Player) (*entities.Player, error)
	FindByID(playerID string) (*entities.Player, error)
}
