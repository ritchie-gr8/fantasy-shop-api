package service

import (
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_adminModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/model"
	_adminRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/repository"
	_playerModel "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/model"
	_playerRepo "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/repository"
)

type googleOAuth2Service struct {
	playerRepo _playerRepo.PlayerRepository
	adminRepo  _adminRepo.AdminRepository
}

func NewGoogleOAuth2Service(
	playerRepo _playerRepo.PlayerRepository,
	adminRepo _adminRepo.AdminRepository,
) OAuth2Service {
	return &googleOAuth2Service{
		playerRepo: playerRepo,
		adminRepo:  adminRepo,
	}
}

func (s *googleOAuth2Service) CreatePlayerAccount(createPlayerReq _playerModel.CreatePlayerReq) error {
	if !s.IsPlayerExistsAndValid(createPlayerReq.ID) {
		player := &entities.Player{
			ID:     createPlayerReq.ID,
			Name:   createPlayerReq.Name,
			Email:  createPlayerReq.Email,
			Avatar: createPlayerReq.Avatar,
		}

		if _, err := s.playerRepo.Create(player); err != nil {
			return err
		}
	}

	return nil
}

func (s *googleOAuth2Service) CreateAdminAccout(createAdminReq _adminModel.CreateAdminReq) error {
	if !s.IsAdminExistsAndValid(createAdminReq.ID) {
		admin := &entities.Admin{
			ID:     createAdminReq.ID,
			Name:   createAdminReq.Name,
			Email:  createAdminReq.Email,
			Avatar: createAdminReq.Avatar,
		}

		if _, err := s.adminRepo.Create(admin); err != nil {
			return err
		}
	}

	return nil
}

func (s *googleOAuth2Service) IsPlayerExistsAndValid(playerID string) bool {
	player, err := s.playerRepo.FindByID(playerID)
	if err != nil {
		return false
	}

	return player != nil
}

func (s *googleOAuth2Service) IsAdminExistsAndValid(adminID string) bool {
	admin, err := s.adminRepo.FindByID(adminID)
	if err != nil {
		return false
	}

	return admin != nil
}
