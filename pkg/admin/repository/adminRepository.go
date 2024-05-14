package repository

import "github.com/ritchie-gr8/fantasy-shop-api/entities"

type AdminRepository interface {
	Create(admin *entities.Admin) (*entities.Admin, error)
	FindByID(adminID string) (*entities.Admin, error)
}
