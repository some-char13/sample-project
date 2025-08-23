package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"
	NewItem "sample_project/internal/service"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer ctx.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startSrvID := 1
	startResID := 1

	commonCh := make(chan repository.SrvID)

	go NewItem.ProcessItems(commonCh, ctx)
	go NewItem.LogItems(ctx)
	go Shutdown(cancel)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Выходим")
			return
		case <-ticker.C:
			s := service.NewService("озон", "https://ozon.ru", 60)
			s.SetServiceID(startSrvID)

			c := check.NewResult(startSrvID, 200, 32)
			c.SetID(startResID)

			select {
			case commonCh <- s:
			case <-ctx.Done():
				return
			}

			select {
			case commonCh <- c:
			case <-ctx.Done():
				return
			}

			startSrvID++
			startResID++
		}

	}
}

func Shutdown(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs

	fmt.Println("Получен сигнал", sig)
	cancel()

}
