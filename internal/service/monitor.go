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

// type MonitorService struct {
// 	repo   repository.Repository
// 	checks map[int]*serviceCheck
// 	mu     *sync.RWMutex
// 	ctx    context.Context
// 	cancel context.CancelFunc
// 	wg     *sync.WaitGroup
// }

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

// func NewMonitorService(repo repository.Repository) *MonitorService {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	return &MonitorService{
// 		repo:   repo,
// 		checks: make(map[int]*serviceCheck),
// 		mu:     &sync.RWMutex{},
// 		ctx:    ctx,
// 		cancel: cancel,
// 		wg:     &sync.WaitGroup{},
// 	}
// }

func (m *MonitorService) Start() error {
	ctx := context.Background()

	services, err := m.repo.GetServices(ctx)
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
	defer m.mu.Unlock()

	for _, check := range m.checks {
		check.ticker.Stop()
		check.cancel()
	}
	m.checks = make(map[int]*serviceCheck)

	m.wg.Wait()
	log.Println("Monitor service stopped")
}

func (m *MonitorService) StartMonitoring(svc *service.Service) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, exists := m.checks[svc.ID]; exists {
		existing.ticker.Stop()
		existing.cancel()
		delete(m.checks, svc.ID)
	}

	ctx, cancel := context.WithCancel(m.ctx)
	ticker := time.NewTicker(time.Duration(svc.Interval) * time.Second)

	check := &serviceCheck{
		service: svc,
		ticker:  ticker,
		cancel:  cancel,
	}
	m.checks[svc.ID] = check

	m.wg.Add(1)
	go m.monitorService(ctx, svc, ticker)
}

func (m *MonitorService) StopMonitoring(serviceID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if check, exists := m.checks[serviceID]; exists {
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

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", svc.URL, nil)
	if err != nil {
		m.saveCheckResult(ctx, svc.ID, 0, int(time.Since(start).Milliseconds()))
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		m.saveCheckResult(ctx, svc.ID, 0, int(responseTime))
		return
	}
	defer resp.Body.Close()

	m.saveCheckResult(ctx, svc.ID, resp.StatusCode, int(responseTime))
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

	status := make(map[int]bool)
	for id := range m.checks {
		status[id] = true
	}
	return status
}
