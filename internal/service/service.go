package service

import (
	"fmt"
	"sample_project/internal/repository"
	"time"
)

func ProcessItems(c chan repository.SrvID) {
	for i := range c {
		repository.AddItem(i)
	}
}

var lenSrv int
var lenRes int

func LogItems() {
	logTicker := time.NewTicker(200 * time.Millisecond)
	defer logTicker.Stop()

	for range logTicker.C {

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
