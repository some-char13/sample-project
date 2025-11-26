package service

import (
	"errors"
	"regexp"
	"time"

	"sample_project/internal/model/check"
)

type Service struct {
	ID       int       `json:"id"`
	Name     string    `json:"name" binding:"required,min=3"`
	URL      string    `json:"url" binding:"required"`
	Interval int       `json:"interval" binding:"required,min=10"`
	Created  time.Time `json:"created"`
}

type Request struct {
	Name     string `json:"name" binding:"required,min=3"`
	URL      string `json:"url" binding:"required"`
	Interval int    `json:"interval" binding:"required,min=10"`
}

type Status struct {
	Service   Service              `json:"service"`
	LastCheck *check.ResultRequest `json:"lastCheck,omitempty"`
}

func NewService(name, url string, interval int) *Service {
	return &Service{
		Name:     name,
		URL:      url,
		Interval: interval,
		Created:  time.Now().UTC(),
	}
}

func (s *Request) Validate() error {
	matched, _ := regexp.MatchString(`^https?://`, s.URL)
	if !matched {
		return errors.New("URL must start with http:// or https://")
	}

	if s.Interval < 10 {
		return errors.New("interval must be at least 10 seconds")
	}

	return nil
}
