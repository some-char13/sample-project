package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
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
		ID:       1,
		Name:     "test-service",
		URL:      "https://httpbin.org/status/200",
		Interval: 1,
	}

	mockRepo.On("AddCheckResult", mock.Anything, mock.MatchedBy(func(r *check.Result) bool {
		return r.ServiceID == 1
	})).Return(nil).Maybe()

	monitor.StartMonitoring(svc)

	time.Sleep(2 * time.Second)

	monitor.StopMonitoring(svc.ID)
	monitor.Stop()

	mockRepo.AssertExpectations(t)
}

func TestMonitorService_Start_WithServices(t *testing.T) {
	mockRepo := new(MockRepository)
	monitor := NewMonitorService(mockRepo)

	services := []*service.Service{
		{
			ID:       1,
			Name:     "google",
			URL:      "https://google.com",
			Interval: 1,
		},
		{
			ID:       2,
			Name:     "github",
			URL:      "https://github.com",
			Interval: 1,
		},
	}

	mockRepo.On("GetServices", mock.Anything).Return(services, nil)

	mockRepo.On("AddCheckResult", mock.Anything, mock.Anything).Return(nil).Maybe()

	assert.NoError(t, monitor.Start())

	time.Sleep(2 * time.Second)

	monitor.Stop()

	mockRepo.AssertExpectations(t)
}
