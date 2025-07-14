package service

import "time"

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

func (s *Service) ID() int {
	return s.id
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) URL() string {
	return s.url
}

func (s *Service) GetInterval() int {
	return s.interval
}

func (s *Service) CreatedAt() time.Time {
	return s.created
}

func (s *Service) SetID(id int) {
	s.id = id
}

func (s *Service) SetName(name string) {
	s.name = name
}

func (s *Service) SetURL(new_url string) {
	s.url = new_url
}

func (s *Service) SetInterval(seconds int) {
	s.interval = seconds
}
