package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/custom"
	_itemShopModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/model"
	_itemShopService "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemShop/service"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/validation"
)

type itemShopControllerImpl struct {
	itemShopService _itemShopService.ItemShopService
}

func NewItemShopController(
	itemShopService _itemShopService.ItemShopService,
) ItemShopController {
	return &itemShopControllerImpl{
		itemShopService: itemShopService,
	}
}

func (c *itemShopControllerImpl) GetItems(pctx echo.Context) error {
	itemFilter := new(_itemShopModel.ItemFilter)
	customReq := custom.NewCustomRequest(pctx)

	if err := customReq.Bind(itemFilter); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	items, err := c.itemShopService.GetItems(itemFilter)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, items)
}

func (c *itemShopControllerImpl) BuyItem(pctx echo.Context) error {
	playerID, err := validation.GetUserID(pctx, "playerID")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	buyReq := new(_itemShopModel.BuyItemReq)
	customReq := custom.NewCustomRequest(pctx)
	if err := customReq.Bind(buyReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	buyReq.PlayerID = playerID

	playerCoin, err := c.itemShopService.BuyItem(buyReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, playerCoin)
}

func (c *itemShopControllerImpl) SellItem(pctx echo.Context) error {
	playerID, err := validation.GetUserID(pctx, "playerID")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	sellReq := new(_itemShopModel.SellItemReq)
	customReq := custom.NewCustomRequest(pctx)
	if err := customReq.Bind(sellReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	sellReq.PlayerID = playerID

	playerCoin, err := c.itemShopService.SellItem(sellReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, playerCoin)
}
