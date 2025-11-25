package service

import (
	"testing"
	"time"

	"sample_project/internal/model/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMonitorService_StartStop(t *testing.T) {
	mockRepo := new(MockRepository)
	monitor := NewMonitorService(mockRepo)

	mockRepo.On("GetServices", mock.Anything).Return([]*service.Service{}, nil)

	assert.NoError(t, monitor.Start())
	monitor.Stop()

	mockRepo.AssertExpectations(t)
}

func TestMonitorService_StartMonitoring(t *testing.T) {
	mockRepo := new(MockRepository)
	monitor := NewMonitorService(mockRepo)

	svc := &service.Service{
		Id:       1,
		Name:     "lenta",
		Url:      "https://lenta.ru",
		Interval: 3,
	}

	mockRepo.On("AddCheckResult", mock.Anything, mock.AnythingOfType("*check.Result")).Return(nil)

	monitor.StartMonitoring(svc)

	time.Sleep(5 * time.Second)

	monitor.StopMonitoring(svc.Id)
	monitor.Stop()

	mockRepo.AssertExpectations(t)
}

func TestMonitorService_Start_WithServices(t *testing.T) {
	mockRepo := new(MockRepository)
	monitor := NewMonitorService(mockRepo)

	services := []*service.Service{
		{
			Id:       1,
			Name:     "test 1",
			Url:      "https://lenta.ru",
			Interval: 2,
		},
		{
			Id:       2,
			Name:     "test 2",
			Url:      "https://github.com",
			Interval: 2,
		},
	}

	mockRepo.On("GetServices", mock.Anything).Return(services, nil)
	mockRepo.On("AddCheckResult", mock.Anything, mock.AnythingOfType("*check.Result")).Return(nil).Maybe()

	assert.NoError(t, monitor.Start())

	time.Sleep(3 * time.Second)

	monitor.Stop()

	mockRepo.AssertExpectations(t)
}
