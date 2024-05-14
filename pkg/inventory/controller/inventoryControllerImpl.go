package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/custom"
	_inventoryService "github.com/ritchie-gr8/fantasy-shop-api/pkg/inventory/service"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/validation"
)

type inventoryControllerImpl struct {
	inventoryService _inventoryService.InventoryService
	logger           echo.Logger
}

func NewInventoryControllerImpl(
	inventoryService _inventoryService.InventoryService,
	logger echo.Logger,
) InventoryController {
	return &inventoryControllerImpl{
		inventoryService: inventoryService,
		logger:           logger,
	}
}

func (c *inventoryControllerImpl) GetInventory(pctx echo.Context) error {
	playerID, err := validation.GetUserID(pctx, "playerID")
	if err != nil {
		c.logger.Errorf("falied to get player id: %s", err.Error())
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	inventory, err := c.inventoryService.GetInventory(playerID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, inventory)
}
