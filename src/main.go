package main

import (
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"
	NewItem "sample_project/internal/service"
	"time"
)

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	startSrvID := 1
	startResID := 1

	commonCh := make(chan repository.SrvID)

	go NewItem.ProcessItems(commonCh)
	go NewItem.LogItems()

	for range ticker.C {

		s := service.NewService("озон", "https://ozon.ru", 60)
		s.SetServiceID(startSrvID)
		commonCh <- s

		c := check.NewResult(startSrvID, 200, 32)
		c.SetID(startResID)
		commonCh <- c

		startSrvID++
		startResID++
	}

}
