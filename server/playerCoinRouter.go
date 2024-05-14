package server

import (
	_playerCoinController "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/controller"
	_playerCoinRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/repository"
	_playerCoinService "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/service"
)

func (s *Server) initPlayerCoinRouter(m *authorizeMiddleware) {
	router := s.app.Group("/v1/player-coin")

	playerCoinRepo := _playerCoinRepo.NewPlayerCoinRepositoryImpl(s.app.Logger, s.db)
	playerCoinService := _playerCoinService.NewPlayerCoinServiceImpl(playerCoinRepo)
	playerCoinController := _playerCoinController.NewPlayerCoinController(playerCoinService)

	router.POST("", playerCoinController.AddCoin, m.AuthorizePlayer)
	router.GET("", playerCoinController.ShowPlayerCoin, m.AuthorizePlayer)
}
