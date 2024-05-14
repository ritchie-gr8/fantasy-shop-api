package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/config"
	_adminModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/model"
	"github.com/ritchie-gr8/fantasy-shop-api/pkg/custom"
	_oauth2Exception "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/exception"
	_oauth2Model "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/model"
	_oauth2Service "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/service"
	_playerModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/model"
	"golang.org/x/oauth2"
)

type googleOAuth2Controller struct {
	oauth2Service _oauth2Service.OAuth2Service
	oauth2Conf    *config.OAuth2
	logger        echo.Logger
}

var (
	playerGoogleOAuth2 *oauth2.Config
	adminGoogleOAuth2  *oauth2.Config
	once               sync.Once

	accessTokenCookieName  = "act"
	refreshTokenCookieName = "rft"
	stateCookieName        = "state"

	letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func NewGoogleOAuth2Controller(
	oauth2Service _oauth2Service.OAuth2Service,
	oauth2Conf *config.OAuth2,
	logger echo.Logger,
) OAuth2Controller {
	once.Do(func() {
		setGoogleOAuth2Config(oauth2Conf)
	})

	return &googleOAuth2Controller{
		oauth2Service: oauth2Service,
		oauth2Conf:    oauth2Conf,
		logger:        logger,
	}
}

func setGoogleOAuth2Config(oauth2Conf *config.OAuth2) {
	playerGoogleOAuth2 = &oauth2.Config{
		ClientID:     oauth2Conf.ClientID,
		ClientSecret: oauth2Conf.ClientSecret,
		RedirectURL:  oauth2Conf.PlayerRedirectUrl,
		Scopes:       oauth2Conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Conf.Endpoints.AuthUrl,
			TokenURL:      oauth2Conf.Endpoints.TokenUrl,
			DeviceAuthURL: oauth2Conf.Endpoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}

	adminGoogleOAuth2 = &oauth2.Config{
		ClientID:     oauth2Conf.ClientID,
		ClientSecret: oauth2Conf.ClientSecret,
		RedirectURL:  oauth2Conf.AdminRedirectUrl,
		Scopes:       oauth2Conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Conf.Endpoints.AuthUrl,
			TokenURL:      oauth2Conf.Endpoints.TokenUrl,
			DeviceAuthURL: oauth2Conf.Endpoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}
}

func (c *googleOAuth2Controller) PlayerLogin(pctx echo.Context) error {
	state := c.randomState()
	c.setCookie(pctx, stateCookieName, state)

	return pctx.Redirect(http.StatusFound, playerGoogleOAuth2.AuthCodeURL(state))
}

func (c *googleOAuth2Controller) AdminLogin(pctx echo.Context) error {
	state := c.randomState()
	c.setCookie(pctx, stateCookieName, state)

	return pctx.Redirect(http.StatusFound, adminGoogleOAuth2.AuthCodeURL(state))
}

func (c *googleOAuth2Controller) PlayerLoginCallback(pctx echo.Context) error {
	ctx := context.Background()

	if err := retry.Do(func() error {
		return c.callbackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("failed to validate callback: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	token, err := playerGoogleOAuth2.Exchange(ctx, pctx.QueryParam("code"))
	if err != nil {
		c.logger.Errorf("failed to exchange token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	client := playerGoogleOAuth2.Client(ctx, token)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Errorf("failed to get user info token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	createPlayerReq := &_playerModel.CreatePlayerReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Email,
		Avatar: userInfo.Picture,
	}

	if err := c.oauth2Service.CreatePlayerAccount(*createPlayerReq); err != nil {
		c.logger.Errorf("failed to get user info token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.OAuth2Processing{})
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, token.RefreshToken)

	return pctx.JSON(http.StatusOK, &_oauth2Model.LoginResponse{
		Message: "Login successfully",
	})
}

func (c *googleOAuth2Controller) AdminLoginCallback(pctx echo.Context) error {
	ctx := context.Background()

	if err := retry.Do(func() error {
		return c.callbackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("failed to validate callback: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	token, err := adminGoogleOAuth2.Exchange(ctx, pctx.QueryParam("code"))
	if err != nil {
		c.logger.Errorf("failed to exchange token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	client := adminGoogleOAuth2.Client(ctx, token)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Errorf("failed to get user info token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	createAdminReq := &_adminModel.CreateAdminReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Email,
		Avatar: userInfo.Picture,
	}

	if err := c.oauth2Service.CreateAdminAccout(*createAdminReq); err != nil {
		c.logger.Errorf("failed to get user info token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.OAuth2Processing{})
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, token.RefreshToken)

	return pctx.JSON(http.StatusOK, &_oauth2Model.LoginResponse{
		Message: "Login successfully",
	})
}

func (c *googleOAuth2Controller) Logout(pctx echo.Context) error {
	accessToken, err := pctx.Cookie(accessTokenCookieName)
	if err != nil {
		c.logger.Errorf("failed to revoke token: %s", err.Error())
		return custom.Error(pctx, http.StatusBadRequest, &_oauth2Exception.Logout{})
	}

	if err := c.revokeToken(accessToken.Value); err != nil {
		c.logger.Errorf("failed to revoke token: %s", err.Error())
		return custom.Error(pctx, http.StatusInternalServerError, &_oauth2Exception.Logout{})
	}
	c.clearSameSiteCookie(pctx, accessTokenCookieName)
	c.clearSameSiteCookie(pctx, refreshTokenCookieName)

	return pctx.JSON(http.StatusOK, &_oauth2Model.LogoutResponse{
		Message: "Logout successfully",
	})
}

func (c *googleOAuth2Controller) revokeToken(accessToken string) error {
	revokeURL := fmt.Sprintf("%s?token=%s", c.oauth2Conf.RevokeUrl, accessToken)

	resp, err := http.Post(revokeURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		c.logger.Errorf("failed to revoke token: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *googleOAuth2Controller) setCookie(pctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		// Secure: , Set this in prod
	}

	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) clearCookie(pctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}

	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) setSameSiteCookie(pctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// Secure: , Set this in prod
	}

	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) clearSameSiteCookie(pctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}

	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) getUserInfo(client *http.Client) (*_oauth2Model.UserInfo, error) {
	resp, err := client.Get(c.oauth2Conf.UserInfoUrl)
	if err != nil {
		c.logger.Errorf("failed to get user info: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	userInfoInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("failed to read user info: %s", err.Error())
		return nil, err
	}

	userInfo := new(_oauth2Model.UserInfo)
	if err := json.Unmarshal(userInfoInBytes, &userInfo); err != nil {
		c.logger.Errorf("failed to unmarshal user info: %s", err.Error())
		return nil, err
	}

	return userInfo, nil
}

func (c *googleOAuth2Controller) callbackValidating(pctx echo.Context) error {
	state := pctx.QueryParam("state")
	stateFromCookie, err := pctx.Cookie(stateCookieName)
	if err != nil {
		c.logger.Errorf("failed to get state from cookie: %s", err.Error())
		return &_oauth2Exception.Unauthorized{}
	}

	if state != stateFromCookie.Value {
		c.logger.Errorf("invalid state", err)
		return &_oauth2Exception.Unauthorized{}
	}

	c.clearCookie(pctx, stateCookieName)

	return nil
}

func (c *googleOAuth2Controller) randomState() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
