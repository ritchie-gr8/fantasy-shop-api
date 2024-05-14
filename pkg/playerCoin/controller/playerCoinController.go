package controller

import "github.com/labstack/echo/v4"

type PlayerCoinController interface {
	AddCoin(pctx echo.Context) error
	ShowPlayerCoin(pctx echo.Context) error
}
