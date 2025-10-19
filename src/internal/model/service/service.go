package service

import (
	"fmt"
	"time"
)

// type Service struct {
// 	Id       int       `json:"id", required`
// 	Name     string    `json:"name"`
// 	Url      string    `json:"url"`
// 	Interval int       `json:"interval"`
// 	Created  time.Time `json:"created"`
// }

// type Service struct {
// 	Id       int       `json:"id" validate:"required"`
// 	Name     string    `json:"name" validate:"required,min=3"`
// 	Url      string    `json:"url" validate:"required"`
// 	Interval int       `json:"interval" validate:"required"`
// 	Created  time.Time `json:"created"`
// }

type Service struct {
	Id       int       `json:"id" binding:"required"`
	Name     string    `json:"name" binding:"required,min=3"`
	Url      string    `json:"url" binding:"required"`
	Interval int       `json:"interval" binding:"required"`
	Created  time.Time `json:"created"`
}

func NewService(id int, name, url string, interval int) *Service {
	return &Service{
		Id:       id,
		Name:     name,
		Url:      url,
		Interval: interval,
		Created:  time.Now().UTC(),
	}
}

func (s *Service) String() string {
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf(
		"id: %d, name: %s, url: %s, interval: %d, created: %v", s.Id, s.Name, s.Url, s.Interval, s.Created,
	)
}

func (s *Service) UnformString() string {
	if s == nil {
		return "nil"
	}
	// return fmt.Sprintf("%d %s %s %d %s",
	// 	s.Id, s.Name, s.Url, s.Interval, s.Created)
	// return fmt.Sprint(
	// 	s.Id, s.Name, s.Url, s.Interval, s.Created)
	return fmt.Sprintf("%d,%s,%s,%d,%s",
		s.Id, s.Name, s.Url, s.Interval, s.Created)
}

// func (s *Service) GetServiceID() int {
// 	return s.id
// }

// func (s *Service) GetServiceName() string {
// 	return s.name
// }

// func (s *Service) GetServiceURL() string {
// 	return s.url
// }

// func (s *Service) GetServiceInterval() int {
// 	return s.interval
// }

// func (s *Service) ServiceCreatedAt() time.Time {
// 	return s.created
// }

// func (s *Service) SetServiceID(id int) {
// 	s.id = id
// }

// func (s *Service) SetServiceName(name string) {
// 	s.name = name
// }

// func (s *Service) SetServiceURL(new_url string) {
// 	s.url = new_url
// }

// func (s *Service) SetServiceInterval(seconds int) {
// 	s.interval = seconds
// }
