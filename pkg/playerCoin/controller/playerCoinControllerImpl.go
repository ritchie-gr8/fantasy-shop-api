package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/custom"
	_playerCoinModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/model"
	_playerCoinService "github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/service"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/playerCoin/validation"
)

type playerCoinControllerImple struct {
	playerCoinService _playerCoinService.PlayerCoinService
}

func NewPlayerCoinController(playerCoinService _playerCoinService.PlayerCoinService) PlayerCoinController {
	return &playerCoinControllerImple{
		playerCoinService: playerCoinService,
	}
}

func (c *playerCoinControllerImple) AddCoin(pctx echo.Context) error {
	playerID, err := validation.GetUserID(pctx, "playerID")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	addCoinReq := new(_playerCoinModel.AddCoinReq)

	customReq := custom.NewCustomRequest(pctx)
	if err := customReq.Bind(addCoinReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	addCoinReq.PlayerID = playerID

	coin, err := c.playerCoinService.AddCoin(*addCoinReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusCreated, coin)
}

func (c *playerCoinControllerImple) ShowPlayerCoin(pctx echo.Context) error {
	playerID, err := validation.GetUserID(pctx, "playerID")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	result := c.playerCoinService.ShowPlayerCoin(playerID)
	return pctx.JSON(http.StatusOK, result)
}
