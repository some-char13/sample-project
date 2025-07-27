package main

import (
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	NewItem "sample_project/internal/service"
	"time"
)

func main() {
	//ниже закомментирован 10 урок
	// url := service.NewService("какой-то сайт", "https://ozon.ru", 30)

	// url.SetName("abc")

	// fmt.Printf("Создан сервис %s (url: %s) c ID %d\n", url.Name(), url.URL(), url.ID())

	// checkResult := check.NewResult(url.ID(), 200, 30)

	// fmt.Printf("Сервис с ID %d вернул код ответа %d за %d ms", checkResult.GetServiceId(), checkResult.GetResponseCode(), checkResult.GetResponseDuration())

	// тут начинается 12 урок
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startSrvID := 1
	startResID := 1

	for range ticker.C {

		s := service.NewService("озон", "https://ozon.ru", 60)
		s.SetServiceID(startSrvID)

		c := check.NewResult(startSrvID, 200, 32)
		c.SetID(startResID)

		startSrvID++
		startResID++

		NewItem.GenerateItems(s, c)
	}

}
