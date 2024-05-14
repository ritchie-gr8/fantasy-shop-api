package service

import (
	_adminModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/model"
	_playerModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/model"
)

type OAuth2Service interface {
	CreatePlayerAccount(createPlayerReq _playerModel.CreatePlayerReq) error
	CreateAdminAccout(createAdminReq _adminModel.CreateAdminReq) error
	IsPlayerExistsAndValid(playerID string) bool
	IsAdminExistsAndValid(adminID string) bool
}
