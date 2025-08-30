package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sync"
)

var new_srv []*service.Service
var new_res []*check.Result

var muteSrv sync.RWMutex
var muteRes sync.RWMutex

func AddItem(item any) {
	switch v := item.(type) {
	case *service.Service:
		muteSrv.Lock()
		new_srv = append(new_srv, v)
		muteSrv.Unlock()
		SaveServicesToFile("services.json")
	case *check.Result:
		muteRes.Lock()
		new_res = append(new_res, v)
		muteRes.Unlock()
		SaveResultsToFile("results.json")
	default:
		fmt.Printf("Unknown type: %T\n", item)
	}
}

func GetServices() []*service.Service {
	muteSrv.RLock()
	sliceSrv := append([]*service.Service(nil), new_srv...)
	muteSrv.RUnlock()
	return sliceSrv
}

func GetResults() []*check.Result {
	muteRes.RLock()
	sliceRes := append([]*check.Result(nil), new_res...)
	muteRes.RUnlock()
	return sliceRes
}

func SaveServicesToFile(path string) {
	muteSrv.RLock()
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка открытия файла")
	}

	jsonStr, err := json.MarshalIndent(new_srv, "", "   ")
	if err != nil {
		log.Fatal(err)
	} else {
		file.Write(jsonStr)
		file.Close()
	}

	muteSrv.RUnlock()
}

func SaveResultsToFile(path string) {
	muteRes.RLock()
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка открытия файла")
	}

	jsonStr, err := json.MarshalIndent(new_res, "", "   ")
	if err != nil {
		log.Fatal(err)
	} else {
		file.Write(jsonStr)
		file.Close()
	}
	muteRes.RUnlock()
}

func LoadServicesFromFile(path string) error {
	_, err := os.OpenFile(path, os.O_RDWR, 0644)
	if os.IsNotExist(err) {
		return nil
	}
	muteSrv.Lock()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var load []*service.Service
	if err := json.Unmarshal(data, &load); err != nil {
		fmt.Println("Ошибка обработки JSON:", err)
	}

	new_srv = load
	muteSrv.Unlock()

	return nil
}

func LoadResultsFromFile(path string) error {
	_, err := os.OpenFile(path, os.O_RDWR, 0644)
	if os.IsNotExist(err) {
		return nil
	}
	muteRes.Lock()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var loaded []*check.Result
	if err := json.Unmarshal(data, &loaded); err != nil {
		fmt.Println("Ошибка обработки JSON:", err)
	}

	new_res = loaded
	muteRes.Unlock()

	return nil
}

// type SrvID interface {
// 	GetServiceID() int
// }

// var new_srv []*service.Service
// var new_res []*check.Result

// var muteSrv sync.RWMutex
// var muteRes sync.RWMutex

// func AddItem(item SrvID) {
// 	switch v := item.(type) {
// 	case *service.Service:
// 		muteSrv.Lock()
// 		new_srv = append(new_srv, v)
// 		muteSrv.Unlock()
// 	case *check.Result:
// 		muteRes.Lock()
// 		new_res = append(new_res, v)
// 		muteRes.Unlock()
// 	default:
// 		fmt.Printf("Unknown type: %T", item)
// 	}
// }

// func GetServices() []*service.Service {
// 	muteSrv.RLock()
// 	sliceSrv := append([]*service.Service(nil), new_srv...)
// 	muteSrv.RUnlock()
// 	return sliceSrv
// }

// func GetResults() []*check.Result {
// 	muteRes.RLock()
// 	sliceRes := append([]*check.Result(nil), new_res...)
// 	muteRes.RUnlock()
// 	return sliceRes
// }
