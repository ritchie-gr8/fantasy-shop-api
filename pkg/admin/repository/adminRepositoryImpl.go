package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/ritchie-gr8/fantasy-shop-api/databases"
	"github.com/ritchie-gr8/fantasy-shop-api/entities"
	_adminException "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/exception"
)

type adminRepositoryImpl struct {
	logger echo.Logger
	db     databases.Database
}

func NewAdminRepositoryImpl(logger echo.Logger, db databases.Database) AdminRepository {
	return &adminRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (r *adminRepositoryImpl) Create(admin *entities.Admin) (*entities.Admin, error) {
	result := new(entities.Admin)
	if err := r.db.Connect().Create(admin).Scan(result).Error; err != nil {
		r.logger.Errorf("failed to create admin: %s", err.Error())
		return nil, &_adminException.CreateAdmin{AdminID: admin.ID}
	}

	return result, nil
}

func (r *adminRepositoryImpl) FindByID(adminID string) (*entities.Admin, error) {
	r.logger.Debug(adminID)
	admin := new(entities.Admin)
	if err := r.db.Connect().Where("id = ?", adminID).First(admin).Error; err != nil {
		r.logger.Errorf("failed to find admin by ID: %s", err.Error())
		return nil, &_adminException.AdminNotFound{AdminID: adminID}
	}

	return admin, nil
}
