package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"

	"github.com/ritchie-gr8/fantasy-shop-api/pkg/custom"
	_oauth2Exception "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/exception"
)

func (c *googleOAuth2Controller) AuthorizePlayer(pctx echo.Context, next echo.HandlerFunc) error {
	ctx := context.Background()

	tokenSource, err := c.getTokenSource(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	if !tokenSource.Valid() {
		tokenSource, err = c.refreshUserToken(pctx, tokenSource, "player")
		if err != nil {
			return custom.Error(pctx, http.StatusUnauthorized, err)
		}
	}

	client := playerGoogleOAuth2.Client(ctx, tokenSource)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	if !c.oauth2Service.IsPlayerExistsAndValid(userInfo.ID) {
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	pctx.Set("playerID", userInfo.ID)

	return next(pctx)
}

func (c *googleOAuth2Controller) AuthorizeAdmin(pctx echo.Context, next echo.HandlerFunc) error {
	ctx := context.Background()

	tokenSource, err := c.getTokenSource(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	if !tokenSource.Valid() {
		tokenSource, err = c.refreshUserToken(pctx, tokenSource, "admin")
		if err != nil {
			return custom.Error(pctx, http.StatusUnauthorized, err)
		}
	}

	client := adminGoogleOAuth2.Client(ctx, tokenSource)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	if !c.oauth2Service.IsAdminExistsAndValid(userInfo.ID) {
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	pctx.Set("adminID", userInfo.ID)

	return next(pctx)
}

func (c *googleOAuth2Controller) refreshUserToken(pctx echo.Context, token *oauth2.Token, userType string) (*oauth2.Token, error) {
	ctx := context.Background()

	var updatedToken *oauth2.Token
	var err error

	if userType == "player" {
		updatedToken, err = playerGoogleOAuth2.TokenSource(ctx, token).Token()
	} else if userType == "admin" {

		updatedToken, err = adminGoogleOAuth2.TokenSource(ctx, token).Token()
	}
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, updatedToken.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, updatedToken.RefreshToken)

	return updatedToken, nil
}

func (c *googleOAuth2Controller) getTokenSource(pctx echo.Context) (*oauth2.Token, error) {
	accessToken, err := pctx.Cookie(accessTokenCookieName)
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	refreshToken, err := pctx.Cookie(refreshTokenCookieName)
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	return &oauth2.Token{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}
