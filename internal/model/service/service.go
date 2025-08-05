package service

import (
	"fmt"
	"time"
)

type Service struct {
	id       int
	name     string
	url      string
	interval int
	created  time.Time
}

func NewService(name, new_url string, interval int) *Service {
	return &Service{
		name:     name,
		url:      new_url,
		interval: interval,
		created:  time.Now().UTC(),
	}
}

func (s *Service) String() string {
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf(
		"id: %d, name: %s, url: %s, interval: %d, created: %v", s.id, s.name, s.url, s.interval, s.created,
	)
}

func (s *Service) GetServiceID() int {
	return s.id
}

func (s *Service) GetServiceName() string {
	return s.name
}

func (s *Service) GetServiceURL() string {
	return s.url
}

func (s *Service) GetServiceInterval() int {
	return s.interval
}

func (s *Service) ServiceCreatedAt() time.Time {
	return s.created
}

func (s *Service) SetServiceID(id int) {
	s.id = id
}

func (s *Service) SetServiceName(name string) {
	s.name = name
}

func (s *Service) SetServiceURL(new_url string) {
	s.url = new_url
}

func (s *Service) SetServiceInterval(seconds int) {
	s.interval = seconds
}
