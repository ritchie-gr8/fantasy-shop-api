package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/custom"
	_itemManagingModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/model"
	_itemManagingService "github.com/ritchie-gr8/fantasy-shop-api/pkg/itemManaging/service"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/validation"
)

type itemManagingControllerImpl struct {
	itemManagingService _itemManagingService.ItemManagingService
}

func NewItemManagingController(
	itemManagingService _itemManagingService.ItemManagingService,
) ItemManagingController {
	return &itemManagingControllerImpl{
		itemManagingService: itemManagingService,
	}
}

func (c *itemManagingControllerImpl) CreateItem(pctx echo.Context) error {
	adminID, err := validation.GetUserID(pctx, "adminID")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	req := new(_itemManagingModel.CreateItemReq)

	customReq := custom.NewCustomRequest(pctx)
	if err := customReq.Bind(req); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	req.AdminID = adminID

	item, err := c.itemManagingService.CreateItem(req)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusCreated, item)
}

func (c *itemManagingControllerImpl) EditItem(pctx echo.Context) error {
	itemID, err := c.getItemID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	editItemReq := new(_itemManagingModel.EditItemReq)
	customReq := custom.NewCustomRequest(pctx)
	if err := customReq.Bind(editItemReq); err != nil {
		return custom.Error(pctx, http.StatusBadGateway, err)
	}

	item, err := c.itemManagingService.EditItem(itemID, editItemReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, item)
}

func (c *itemManagingControllerImpl) ArchiveItem(pctx echo.Context) error {
	itemID, err := c.getItemID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	if err := c.itemManagingService.ArchiveItem(itemID); err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.NoContent(http.StatusNoContent)
}

func (c *itemManagingControllerImpl) getItemID(pctx echo.Context) (uint64, error) {
	rawItemID := pctx.Param("itemID")
	itemID, err := strconv.ParseUint(rawItemID, 10, 64)
	if err != nil {
		return 0, err
	}

	return itemID, nil
}
