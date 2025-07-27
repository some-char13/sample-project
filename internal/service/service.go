package service

import (
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"
)

func GenerateItems(s *service.Service, c *check.Result) {

	repository.AddItem(s)

	repository.AddItem(c)
}
