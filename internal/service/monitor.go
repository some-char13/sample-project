package service

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"
)

type MonitorService struct {
	repo   repository.Repository
	checks map[int]*serviceCheck
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type serviceCheck struct {
	service *service.Service
	ticker  *time.Ticker
	cancel  context.CancelFunc
}

func NewMonitorService(repo repository.Repository) *MonitorService {
	ctx, cancel := context.WithCancel(context.Background())
	return &MonitorService{
		repo:   repo,
		checks: make(map[int]*serviceCheck),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *MonitorService) Start() error {
	services, err := m.repo.GetServices(context.Background())
	if err != nil {
		return err
	}

	for _, svc := range services {
		m.StartMonitoring(svc)
	}

	log.Printf("Started monitoring for %d services", len(services))
	return nil
}

func (m *MonitorService) Stop() {
	m.cancel()

	m.mu.Lock()
	for _, check := range m.checks {
		check.ticker.Stop()
		check.cancel()
	}
	m.checks = make(map[int]*serviceCheck)
	m.mu.Unlock()

	m.wg.Wait()
	log.Println("Monitor service stopped")
}

func (m *MonitorService) StartMonitoring(svc *service.Service) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if old, ok := m.checks[svc.ID]; ok {
		old.ticker.Stop()
		old.cancel()
		delete(m.checks, svc.ID)
	}

	childCtx, cancel := context.WithCancel(m.ctx)
	ticker := time.NewTicker(time.Duration(svc.Interval) * time.Second)

	m.checks[svc.ID] = &serviceCheck{
		service: svc,
		ticker:  ticker,
		cancel:  cancel,
	}

	m.wg.Add(1)
	go m.monitorService(childCtx, svc, ticker)
}

func (m *MonitorService) StopMonitoring(serviceID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if check, ok := m.checks[serviceID]; ok {
		check.ticker.Stop()
		check.cancel()
		delete(m.checks, serviceID)
	}
}

func (m *MonitorService) monitorService(ctx context.Context, svc *service.Service, ticker *time.Ticker) {
	defer m.wg.Done()

	m.performCheck(ctx, svc)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.performCheck(ctx, svc)
		}
	}
}

func (m *MonitorService) performCheck(ctx context.Context, svc *service.Service) {
	start := time.Now()

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, svc.URL, nil)
	req.Header.Set("User-Agent", "MonitorBot/1.0")

	resp, err := client.Do(req)
	duration := int(time.Since(start).Milliseconds())

	statusCode := 0
	if err == nil {
		statusCode = resp.StatusCode
		resp.Body.Close()
	}

	m.saveCheckResult(ctx, svc.ID, statusCode, duration)
}

func (m *MonitorService) saveCheckResult(ctx context.Context, serviceID, statusCode, responseTime int) {
	result := check.NewResult(serviceID, statusCode, responseTime)
	if err := m.repo.AddCheckResult(ctx, result); err != nil {
		log.Printf("Failed to save check result for service %d: %v", serviceID, err)
	}
}

func (m *MonitorService) GetMonitoringStatus() map[int]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[int]bool, len(m.checks))
	for id := range m.checks {
		status[id] = true
	}
	return status
}
