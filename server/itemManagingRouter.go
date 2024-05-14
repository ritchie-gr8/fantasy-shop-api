package server

import (
	_itemManagingController "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/controller"
	_itemManagingRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/repository"
	_itemManagingService "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/service"
	_itemShopRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/repository"
)

func (s *Server) initItemManagingRouter(m *authorizeMiddleware) {
	router := s.app.Group("/v1/item-managing")

	itemShopRepo := _itemShopRepo.NewItemShopRepositoryImpl(s.app.Logger, s.db)
	itemManagingRepo := _itemManagingRepo.NewItemManagingRepositoryImpl(s.app.Logger, s.db)
	itemManagingService := _itemManagingService.NewItemShopServiceImpl(itemManagingRepo, itemShopRepo)
	itemManagingController := _itemManagingController.NewItemManagingController(itemManagingService)

	router.POST("", itemManagingController.CreateItem, m.AuthorizeAdmin)
	router.PATCH("/:itemID", itemManagingController.EditItem, m.AuthorizeAdmin)
	router.DELETE("/:itemID", itemManagingController.ArchiveItem, m.AuthorizeAdmin)
}
