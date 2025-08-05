package repository

import (
	"fmt"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sync"
	"time"
)

type SrvID interface {
	GetServiceID() int
}

var new_srv []*service.Service
var new_res []*check.Result

var muteSrv sync.RWMutex
var muteRes sync.RWMutex

func AddItem(item SrvID) {
	switch v := item.(type) {
	case *service.Service:
		muteSrv.Lock()
		new_srv = append(new_srv, v)
		muteSrv.Unlock()
	case *check.Result:
		muteRes.Lock()
		new_res = append(new_res, v)
		muteRes.Unlock()
	default:
		fmt.Printf("Unknown type: %T", item)
	}
}

func LogItems() {
	var len_srv int
	var len_res int

	logTicker := time.NewTicker(200 * time.Millisecond)
	defer logTicker.Stop()
	for range logTicker.C {

		muteSrv.RLock()
		if len_srv < len(new_srv) {
			new_srv_len := len(new_srv)
			print_this := new_srv[len_srv:new_srv_len]
			for _, item := range print_this {
				fmt.Println(item)
			}
		}
		len_srv = len(new_srv)
		muteSrv.RUnlock()

		muteRes.RLock()
		if len_res < len(new_res) {
			new_res_len := len(new_res)
			print_this := new_res[len_res:new_res_len]
			for _, item := range print_this {
				fmt.Println(item)
			}
		}
		len_res = len(new_res)
		muteRes.RUnlock()
	}
}
