package server

import (
	_inventoryRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/repository"
	_itemShopController "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/controller"
	_itemShopRepository "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
	_itemShopService "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/service"
	_playerCoinRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/repository"
)

func (s *Server) initItemShopRouter(m *authorizeMiddleware) {
	router := s.app.Group("/v1/item-shop")

	inventoryRepo := _inventoryRepo.NewInventoryRepositoryImpl(s.app.Logger, s.db)
	playerCoinRepo := _playerCoinRepo.NewPlayerCoinRepositoryImpl(s.app.Logger, s.db)
	itemShopRepository := _itemShopRepository.NewItemShopRepositoryImpl(s.app.Logger, s.db)
	itemShopService := _itemShopService.NewItemShopServiceImpl(
		itemShopRepository,
		playerCoinRepo,
		inventoryRepo,
		s.app.Logger,
	)
	itemShopController := _itemShopController.NewItemShopController(itemShopService)

	router.GET("", itemShopController.GetItems)
	router.POST("/buy", itemShopController.BuyItem, m.AuthorizePlayer)
	router.POST("/sell", itemShopController.SellItem, m.AuthorizePlayer)
}
