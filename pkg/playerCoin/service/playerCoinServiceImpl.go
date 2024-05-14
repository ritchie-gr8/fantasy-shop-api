package service

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	_playerCoinRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/repository"
)

type playerCoinServiceImple struct {
	playerCoinRepo _playerCoinRepo.PlayerCoinRepository
}

func NewPlayerCoinServiceImpl(playerCoinRepo _playerCoinRepo.PlayerCoinRepository) PlayerCoinService {
	return &playerCoinServiceImple{
		playerCoinRepo: playerCoinRepo,
	}
}

func (s *playerCoinServiceImple) AddCoin(addCoinReq _playerCoinModel.AddCoinReq) (*_playerCoinModel.PlayerCoin, error) {
	coin := &entities.PlayerCoin{
		PlayerID: addCoinReq.PlayerID,
		Amount:   addCoinReq.Amount,
	}

	result, err := s.playerCoinRepo.AddCoin(nil, coin)
	if err != nil {
		return nil, err
	}

	return result.ToPlayerCoinModel(), nil
}

func (s *playerCoinServiceImple) ShowPlayerCoin(playerID string) *_playerCoinModel.ShowPlayerCoin {
	result, err := s.playerCoinRepo.ShowCoin(playerID)
	if err != nil {
		return &_playerCoinModel.ShowPlayerCoin{
			PlayerID: playerID,
			Coin:     0,
		}
	}

	return result
}
