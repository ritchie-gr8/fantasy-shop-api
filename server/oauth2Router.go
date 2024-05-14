package server

import (
	_adminRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/repository"
	_oauth2Controller "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/controller"
	_oauth2Service "github.com/ritchie-gr8/fantasy-shop-api/pkg/oauth2/service"
	_playerRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/repository"
)

func (s *Server) initOAuth2Router() {
	router := s.app.Group("/v1/oauth2/google")

	playerRepo := _playerRepo.NewPlayerRepositoryImpl(s.app.Logger, s.db)
	_adminRepo := _adminRepo.NewAdminRepositoryImpl(s.app.Logger, s.db)
	oauth2Service := _oauth2Service.NewGoogleOAuth2Service(playerRepo, _adminRepo)
	oauth2Controller := _oauth2Controller.NewGoogleOAuth2Controller(
		oauth2Service,
		s.conf.OAuth2,
		s.app.Logger,
	)

	router.GET("/player/login", oauth2Controller.PlayerLogin)
	router.GET("/admin/login", oauth2Controller.AdminLogin)
	router.GET("/player/login/callback", oauth2Controller.PlayerLoginCallback)
	router.GET("/admin/login/callback", oauth2Controller.AdminLoginCallback)
	router.DELETE("/logout", oauth2Controller.Logout)
}
