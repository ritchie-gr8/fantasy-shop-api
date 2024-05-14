package controller

import "github.com/labstack/echo/v4"

type InventoryController interface {
	GetInventory(pctx echo.Context) error
}
