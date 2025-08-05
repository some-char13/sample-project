package service

import (
	"sample_project/internal/repository"
)

func GenerateItems(c chan repository.SrvID) {
	for i := range c {
		go repository.AddItem(i)
	}
}
