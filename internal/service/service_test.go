package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddService(
	ctx context.Context,
	svc *service.Service,
) (*service.Service, error) {
	args := m.Called(ctx, svc)
	return args.Get(0).(*service.Service), args.Error(1)
}

func (m *MockRepository) GetServiceByName(ctx context.Context, name string) (*service.Service, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Service), args.Error(1)
}

func (m *MockRepository) GetServices(ctx context.Context) ([]*service.Service, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*service.Service), args.Error(1)
}

func (m *MockRepository) GetServiceByID(ctx context.Context, id int) (*service.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Service), args.Error(1)
}

func (m *MockRepository) DeleteService(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) AddCheckResult(ctx context.Context, result *check.Result) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *MockRepository) GetCheckResults(ctx context.Context, serviceID, limit int) ([]*check.Result, error) {
	args := m.Called(ctx, serviceID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*check.Result), args.Error(1)
}

func (m *MockRepository) GetCheckResultsByStatus(
	ctx context.Context,
	serviceID int,
	respCodes []int,
	limit int,
) ([]*check.Result, error) {
	args := m.Called(ctx, serviceID, respCodes, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*check.Result), args.Error(1)
}

func (m *MockRepository) GetLastServiceCheck(ctx context.Context, serviceID int) (*check.ResultRequest, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*check.ResultRequest), args.Error(1)
}

func (m *MockRepository) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateService(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := &service.Request{
		Name:     "Test Service",
		URL:      "https://example.com",
		Interval: 30,
	}

	mockRepo.On("GetServiceByName", mock.Anything, "Test Service").
		Return((*service.Service)(nil), repository.ErrServiceNotFound)

	createdService := &service.Service{
		ID:       1,
		Name:     "Test Service",
		URL:      "https://example.com",
		Interval: 30,
	}
	mockRepo.On("AddService", mock.Anything, mock.Anything).
		Return(createdService, nil)

	result, err := svc.CreateService(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Service", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateService_DuplicateName(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	existing := &service.Service{ID: 1, Name: "Test Service"}
	mockRepo.On("GetServiceByName", mock.Anything, "Test Service").
		Return(existing, nil)

	_, err := svc.CreateService(context.Background(), &service.Request{
		Name: "Test Service", URL: "https://example.com", Interval: 30,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestCreateService_ValidationError(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := &service.Request{
		Name:     "Test",
		URL:      "Qwerty",
		Interval: 30,
	}

	result, err := svc.CreateService(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "URL must start with http:// or https://")
}

func TestGetServices(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	expectedServices := []*service.Service{
		{ID: 1, Name: "Service 1", URL: "https://lenta.ru", Interval: 30},
		{ID: 2, Name: "Service 2", URL: "https://google.com", Interval: 60},
	}

	mockRepo.On("GetServices", mock.Anything).Return(expectedServices, nil)

	services, err := svc.GetServices(context.Background())

	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "Service 1", services[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestGetServiceStatus(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	existingService := &service.Service{
		ID:       1,
		Name:     "Test Service",
		URL:      "https://lenta.ru",
		Interval: 30,
	}

	latestCheck := &check.ResultRequest{
		ServiceID:    1,
		ResponseCode: 200,
		RespDuration: 150,
		TimeChecked:  time.Now(),
	}

	mockRepo.On("GetServiceByID", mock.Anything, 1).
		Return(existingService, nil)
	mockRepo.On("GetLastServiceCheck", mock.Anything, 1).
		Return(latestCheck, nil)

	status, err := svc.GetServiceStatus(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, 1, status.Service.ID)
	assert.Equal(t, 200, status.LastCheck.ResponseCode)
	mockRepo.AssertExpectations(t)
}

func TestDeleteService(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	mockRepo.On("DeleteService", mock.Anything, 1).Return(nil)

	err := svc.DeleteService(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetServiceResults(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	expectedResults := []*check.Result{
		{
			ID:           1,
			ServiceID:    1,
			ResponseCode: 200,
			RespDuration: 150,
			TimeChecked:  time.Now().UTC(),
		},
		{
			ID:           2,
			ServiceID:    1,
			ResponseCode: 200,
			RespDuration: 200,
			TimeChecked:  time.Now().UTC().Add(-time.Minute),
		},
	}

	mockRepo.On("GetCheckResults", mock.Anything, 1, 20).Return(expectedResults, nil)

	results, err := svc.GetServiceResults(context.Background(), 1, 20)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 200, results[0].ResponseCode)
	mockRepo.AssertExpectations(t)
}
