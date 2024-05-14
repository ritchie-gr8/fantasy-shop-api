package server

import (
	_inventoryController "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/controller"
	_inventoryRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/repository"
	_inventoryService "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/service"
	_itemShopRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
)

func (s *Server) initInventoryRouter(m *authorizeMiddleware) {
	router := s.app.Group("/v1/inventory")

	itemShopRepo := _itemShopRepo.NewItemShopRepositoryImpl(s.app.Logger, s.db)
	inventoryRepo := _inventoryRepo.NewInventoryRepositoryImpl(s.app.Logger, s.db)
	inventoryService := _inventoryService.NewInventoryServiceImpl(inventoryRepo, itemShopRepo)
	inventoryController := _inventoryController.NewInventoryControllerImpl(inventoryService, s.app.Logger)

	router.GET("", inventoryController.GetInventory, m.AuthorizePlayer)
}
