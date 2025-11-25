package service

import (
	"errors"
	"regexp"
	"sample_project/internal/model/check"
	"time"
)

type Service struct {
	Id       int       `json:"id"`
	Name     string    `json:"name" binding:"required,min=3"`
	Url      string    `json:"url" binding:"required"`
	Interval int       `json:"interval" binding:"required,min=10"`
	Created  time.Time `json:"created"`
}

type ServiceRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Url      string `json:"url" binding:"required"`
	Interval int    `json:"interval" binding:"required,min=10"`
}

type ServiceStatus struct {
	Service   Service              `json:"service"`
	LastCheck *check.ResultRequest `json:"last_check,omitempty"`
}

func NewService(name, url string, interval int) *Service {
	return &Service{
		Name:     name,
		Url:      url,
		Interval: interval,
		Created:  time.Now().UTC(),
	}
}

func (s *ServiceRequest) Validate() error {
	matched, _ := regexp.MatchString(`^https?://`, s.Url)
	if !matched {
		return errors.New("URL must start with http:// or https://")
	}

	if s.Interval < 10 {
		return errors.New("interval must be at least 10 seconds")
	}

	return nil
}
