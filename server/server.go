package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ritchie-gr8/fantasy-shop-api/config"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"

	_adminRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/repository"
	_oauth2Controller "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/controller"
	_oauth2Service "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/service"
	_playerRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/repository"
)

type Server struct {
	app  *echo.Echo
	db   databases.Database
	conf *config.Config
}

var (
	once   sync.Once
	server *Server
)

func NewServer(conf *config.Config, db databases.Database) *Server {
	app := echo.New()
	app.Logger.SetLevel(log.DEBUG)

	once.Do(func() {
		server = &Server{
			app:  app,
			db:   db,
			conf: conf,
		}
	})

	return server
}

func (s *Server) Start() {

	corsMiddleware := getCORSMiddleware(s.conf.Server.AllowOrigins)
	bodyLimitMiddleware := getBodyLimitMiddleware(s.conf.Server.BodyLimit)
	timeoutMiddleware := getTimeOutMiddleware(s.conf.Server.TimeOut)

	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())
	s.app.Use(corsMiddleware)
	s.app.Use(bodyLimitMiddleware)
	s.app.Use(timeoutMiddleware)

	authorizeMiddleware := s.getAuthorizeMiddleware()

	s.app.GET("/v1/health", s.heathCheck)

	s.initOAuth2Router()
	s.initItemShopRouter(authorizeMiddleware)
	s.initItemManagingRouter(authorizeMiddleware)
	s.initPlayerCoinRouter(authorizeMiddleware)
	s.initInventoryRouter(authorizeMiddleware)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	go s.gracefullyShutdown(quitCh)
	s.listen()
}

func (s *Server) listen() {
	url := fmt.Sprintf(":%d", s.conf.Server.Port)

	if err := s.app.Start(url); err != nil && err != http.ErrServerClosed {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}
}

func (s *Server) gracefullyShutdown(quitCh chan os.Signal) {
	ctx := context.Background()

	<-quitCh
	s.app.Logger.Info("Shutting down server...")

	if err := s.app.Shutdown(ctx); err != nil {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}
}

func (s *Server) heathCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func getTimeOutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request Timeout",
		Timeout:      timeout * time.Second,
	})
}

func getCORSMiddleware(allowOrigins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	})
}

func getBodyLimitMiddleware(bodyLimit string) echo.MiddlewareFunc {
	return middleware.BodyLimit(bodyLimit)
}

func (s *Server) getAuthorizeMiddleware() *authorizeMiddleware {
	playerRepo := _playerRepo.NewPlayerRepositoryImpl(s.app.Logger, s.db)
	_adminRepo := _adminRepo.NewAdminRepositoryImpl(s.app.Logger, s.db)
	oauth2Service := _oauth2Service.NewGoogleOAuth2Service(playerRepo, _adminRepo)
	oauth2Controller := _oauth2Controller.NewGoogleOAuth2Controller(
		oauth2Service,
		s.conf.OAuth2,
		s.app.Logger,
	)

	return &authorizeMiddleware{
		oauth2Controller: oauth2Controller,
		oauth2Conf:       s.conf.OAuth2,
		logger:           s.app.Logger,
	}
}
