package main

import (
	"fmt"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
)

func main() {

	url := service.NewService("какой-то сайт", "https://ozon.ru", 30)

	url.SetName("abc")

	fmt.Printf("Создан сервис %s (url: %s) c ID %d\n", url.Name(), url.URL(), url.ID())

	checkResult := check.NewResult(url.ID(), 200, 30)

	fmt.Printf("Сервис с ID %d вернул код ответа %d за %d ms", checkResult.GetServiceId(), checkResult.GetResponseCode(), checkResult.GetResponseDuration())

}
