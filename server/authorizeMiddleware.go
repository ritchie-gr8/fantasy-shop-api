package server

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/config"
	_oauth2Controller "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/controller"
)

type authorizeMiddleware struct {
	oauth2Controller _oauth2Controller.OAuth2Controller
	oauth2Conf       *config.OAuth2
	logger           echo.Logger
}

func (m *authorizeMiddleware) AuthorizePlayer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(pctx echo.Context) error {
		return m.oauth2Controller.AuthorizePlayer(pctx, next)
	}
}

func (m *authorizeMiddleware) AuthorizeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(pctx echo.Context) error {
		return m.oauth2Controller.AuthorizeAdmin(pctx, next)
	}
}
