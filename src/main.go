package main

import (
	"sample_project/internal/app"
	"sample_project/internal/repository"
)

func main() {
	repository.LoadCsvServices("services.csv")
	repository.LoadCsvResults("results.csv")
	// service.SetCount(len(repository.GetServices()), len(repository.GetResults()))
	//_, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//go Shutdown(cancel)

	application := app.New()

	application.Start()

	application.Stop()
}

// func main() {

// 	repository.LoadServicesFromFile("services.json")
// 	repository.LoadResultsFromFile("results.json")

// 	NewItem.SetCount(len(repository.GetServices()), len(repository.GetResults()))

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	startSrvID := 1
// 	startResID := 1

// 	commonCh := make(chan any)

// 	go NewItem.ProcessItems(commonCh, ctx)
// 	//go NewItem.LogItems(ctx)
// 	go Shutdown(cancel)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("Выходим")
// 			return
// 		case <-ticker.C:
// 			//s := service.NewService("озон", "https://ozon.ru", 60)
// 			s := service.NewService("qaz", "https://ozon.ru", 60)
// 			s.Id = startSrvID

// 			c := check.NewResult(startSrvID, 0, 32)
// 			c.Id = startResID

// 			select {
// 			case commonCh <- s:
// 			case <-ctx.Done():
// 				return
// 			}

// 			select {
// 			case commonCh <- c:
// 			case <-ctx.Done():
// 				return
// 			}

// 			startSrvID++
// 			startResID++
// 		}
// 	}
// }

// func Shutdown(cancel context.CancelFunc) {
// 	sigs := make(chan os.Signal, 1)
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
// 	sig := <-sigs
// 	fmt.Println("Получен сигнал", sig)
// 	cancel()
// }
