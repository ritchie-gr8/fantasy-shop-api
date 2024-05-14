package service

import (
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
)

type PlayerCoinService interface {
	AddCoin(addCoinReq _playerCoinModel.AddCoinReq) (*_playerCoinModel.PlayerCoin, error)
	ShowPlayerCoin(playerID string) *_playerCoinModel.ShowPlayerCoin
}
