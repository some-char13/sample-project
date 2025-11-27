package repository

import (
	"context"
	"errors"

	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
)

var (
	ErrServiceNotFound = errors.New("service not found")
	ErrCheckNotFound   = errors.New("check result not found")
)

type Repository interface {
	AddService(ctx context.Context, s *service.Service) (*service.Service, error)
	GetServices(ctx context.Context) ([]*service.Service, error)
	GetServiceByID(ctx context.Context, id int) (*service.Service, error)
	GetServiceByName(ctx context.Context, name string) (*service.Service, error)
	DeleteService(ctx context.Context, id int) error

	AddCheckResult(ctx context.Context, result *check.Result) error
	GetCheckResults(ctx context.Context, serviceID int, limit int) ([]*check.Result, error)
	GetLastServiceCheck(ctx context.Context, serviceID int) (*check.ResultRequest, error)
	GetCheckResultsByStatus(ctx context.Context, serviceID int, respCodes []int, limit int) ([]*check.Result, error)

	Ping() error
	Close() error
}
