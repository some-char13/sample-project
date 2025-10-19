package service

import (
	"context"
	"fmt"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"
	"time"
)

// func ProcessItems(c chan any, ctx context.Context) {
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

// func ProcessItems(c chan any) {
// 	item := <-c
// 	repository.AddItem(item)
// }

// func ChangeItems(i int, c chan any) {
// 	item := <-c
// 	repository.ChangeItem(i, item)
// }

func ProcessItems(i any) {
	repository.AddItem(i)
}

func ChangeItems(i int, item any) {
	repository.ChangeItem(i, item)
}

func GetServices() []*service.Service {
	srv := repository.GetServices()
	return srv
}

func GetResults() []*check.Result {
	res := repository.GetResults()
	return res
}

func SearchServiceItem(i int) *service.Service {
	srv := repository.SearchItemService(i)
	return srv
}

func SearchResultItem(i int) *check.Result {
	//item := <-c
	res := repository.SearchItemResult(i)
	return res
}

func DeleteItemService(i int) {
	//item := <-c
	repository.DeleteItemService(i)
}

func DeleteItemResult(i int) {
	//item := <-c
	repository.DeleteItemResult(i)
}

var lenSrv int
var lenRes int

func SetCount(serviceCnt, resultCnt int) {
	lenSrv = serviceCnt
	lenRes = resultCnt
}

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

// func ChangeItem(c chan any, ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("Завершаем обработку")
// 			return
// 		case item := <-c:
// 			repository.ChangeItem(item)
// 			fmt.Println("Отправили в AddItem")
// 		}
// 	}
// }

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
