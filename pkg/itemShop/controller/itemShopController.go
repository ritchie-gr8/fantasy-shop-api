package controller

import "github.com/labstack/echo/v4"

type ItemShopController interface {
	GetItems(pctx echo.Context) error
	BuyItem(pctx echo.Context) error
	SellItem(pctx echo.Context) error
}
