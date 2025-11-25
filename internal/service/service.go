package service

import (
	"context"
	"errors"
	"fmt"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateService(ctx context.Context, serviceReq *service.ServiceRequest) (*service.Service, error) {
	if err := serviceReq.Validate(); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetServiceByName(ctx, serviceReq.Name)
	if err != nil {
		if errors.Is(err, repository.ErrServiceNotFound) {
		} else {
			return nil, fmt.Errorf("failed to check service existence: %w", err)
		}
	} else if existing != nil {
		return nil, fmt.Errorf("service with name '%s' already exists", serviceReq.Name)
	}

	svc := service.NewService(serviceReq.Name, serviceReq.Url, serviceReq.Interval)

	createdSvc, err := s.repo.AddService(ctx, svc)
	if err != nil {
		return nil, err
	}

	return createdSvc, nil
}

func (s *Service) GetServices(ctx context.Context) ([]*service.Service, error) {
	return s.repo.GetServices(ctx)
}

func (s *Service) GetServiceByID(ctx context.Context, id int) (*service.Service, error) {
	chk, err := s.repo.GetServiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrServiceNotFound) {
			return nil, fmt.Errorf("service with id %d not found", id)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return chk, nil
}

func (s *Service) DeleteService(ctx context.Context, id int) error {
	return s.repo.DeleteService(ctx, id)
}

func (s *Service) GetServiceResults(ctx context.Context, serviceID int, limit int) ([]*check.Result, error) {
	return s.repo.GetCheckResults(ctx, serviceID, limit)
}

func (s *Service) GetServiceResultsByStatus(ctx context.Context, serviceID int, respCodes []int, limit int) ([]*check.Result, error) {
	return s.repo.GetCheckResultsByStatus(ctx, serviceID, respCodes, limit)
}

func (s *Service) GetServiceStatus(ctx context.Context, serviceID int) (*service.ServiceStatus, error) {
	svc, err := s.repo.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, repository.ErrServiceNotFound) {
			return nil, fmt.Errorf("service with id %d not found", serviceID)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	latest, err := s.repo.GetLastServiceCheck(ctx, serviceID)
	if err != nil {
		if errors.Is(err, repository.ErrCheckNotFound) {
			return nil, fmt.Errorf("check result for id %d not found", serviceID)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	status := &service.ServiceStatus{
		Service:   *svc,
		LastCheck: latest,
	}

	return status, nil
}

func (s *Service) GetAllServicesStatus(ctx context.Context) ([]*service.ServiceStatus, error) {
	services, err := s.repo.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	var statuses []*service.ServiceStatus
	for _, svc := range services {
		status, err := s.GetServiceStatus(ctx, svc.Id)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}
