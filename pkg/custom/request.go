package custom

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	Request interface {
		Bind(obj any) error
	}

	customRequest struct {
		ctx       echo.Context
		validator *validator.Validate
	}
)

var (
	once              sync.Once
	validatorInstance *validator.Validate
)

func NewCustomRequest(req echo.Context) Request {
	once.Do(func() {
		validatorInstance = validator.New()
	})

	return &customRequest{
		ctx:       req,
		validator: validatorInstance,
	}
}

func (r *customRequest) Bind(obj any) error {
	if err := r.ctx.Bind(obj); err != nil {
		return err
	}

	if err := r.validator.Struct(obj); err != nil {
		return err
	}

	return nil
}
