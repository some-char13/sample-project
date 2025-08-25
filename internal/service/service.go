package service

import (
	"context"
	"fmt"
	"sample_project/internal/repository"
	"time"
)

func ProcessItems(c chan any, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Завершаем обработку")
			return
		case item := <-c:
			repository.AddItem(item)
		}
	}
}

var lenSrv int
var lenRes int

func LogItems(ctx context.Context) {
	logTicker := time.NewTicker(200 * time.Millisecond)
	defer logTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Завершаем логирование")
			return
		case <-logTicker.C:
			newSrv := repository.GetServices()
			if lenSrv < len(newSrv) {
				printSrv := newSrv[lenSrv:]
				for _, s := range printSrv {
					fmt.Println(s)
				}
			}
			lenSrv = len(newSrv)

			newRes := repository.GetResults()
			if lenRes < len(newRes) {
				printRes := newRes[lenRes:]
				for _, r := range printRes {
					fmt.Println(r)
				}
			}
			lenRes = len(newRes)
		}
	}
}

// func ProcessItems(c chan repository.SrvID, ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("Завершаем обработку")
// 			return
// 		case item := <-c:
// 			repository.AddItem(item)
// 		}
// 	}
// }

// var lenSrv int
// var lenRes int

// func LogItems(ctx context.Context) {
// 	logTicker := time.NewTicker(200 * time.Millisecond)
// 	defer logTicker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("Завершаем логирование")
// 			return
// 		case <-logTicker.C:

// 			newSrv := repository.GetServices()
// 			if lenSrv < len(newSrv) {
// 				printSrv := newSrv[lenSrv:]
// 				for _, s := range printSrv {
// 					fmt.Println(s)
// 				}
// 			}
// 			lenSrv = len(newSrv)

// 			newRes := repository.GetResults()
// 			if lenRes < len(newRes) {
// 				printRes := newRes[lenRes:]
// 				for _, r := range printRes {
// 					fmt.Println(r)
// 				}
// 			}
// 			lenRes = len(newRes)
// 		}

// 	}
// }
