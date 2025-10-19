package repository

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"strconv"
	"strings"
	"sync"
	"time"
)

var new_srv []*service.Service
var new_res []*check.Result
var Abc *service.Service

var muteSrv sync.RWMutex
var muteRes sync.RWMutex

func AddItem(item any) {
	//validate := validator.New()
	switch v := item.(type) {
	case *service.Service:
		// err := validate.Struct(v)
		// if err != nil {
		// 	for _, err := range err.(validator.ValidationErrors) {
		// 		fmt.Printf("Ошибка в поле %s: %s\n", err.Field(), err.Tag())
		// 		//muteSrv.Unlock()
		// 		return
		// 	}
		// }
		muteSrv.Lock()
		//fmt.Println("до append", new_srv)
		new_srv = append(new_srv, v)
		//fmt.Println("после append", new_srv)
		muteSrv.Unlock()
		if err := SaveServicesToFileCsv2("services.csv", new_srv); err != nil {
			fmt.Println("Ошибка записи:", err)
		}
		//fmt.Println("Записали в SaveServicesToFileCsv2")
		//SaveServicesToFile("services.json")
	case *check.Result:
		muteRes.Lock()
		new_res = append(new_res, v)
		muteRes.Unlock()
		if err := SaveResultsToFileCsv2("results.csv", new_res); err != nil {
			fmt.Println("Ошибка записи:", err)
		}
		//		SaveResultsToFile("results.json")
	default:
		fmt.Printf("Unknown type: %T\n", item)
	}
}

func ChangeItem(i int, item any) {
	//validate := validator.New()
	switch v := item.(type) {
	case *service.Service:
		// err := validate.Struct(v)
		// if err != nil {
		// 	for _, err := range err.(validator.ValidationErrors) {
		// 		fmt.Printf("Ошибка в поле %s: %s\n", err.Field(), err.Tag())
		// 		//muteSrv.Unlock()
		// 		return
		// 	}
		// }

		muteSrv.Lock()
		for i, val := range new_srv {
			if val.Id == v.Id {

				sl := new_srv[i]
				sl.Name = v.Name
				sl.Url = v.Url
				sl.Interval = v.Interval

			}
		}
		muteSrv.Unlock()
		if err := SaveServicesToFileCsv2("services.csv", new_srv); err != nil {
			fmt.Println("Ошибка записи:", err)
		}
	case *check.Result:
		muteRes.Lock()
		for i, val := range new_res {
			if val.Id == v.Id {

				sl := new_res[i]
				sl.ServiceId = v.ServiceId
				sl.ResponseCode = v.ResponseCode
				sl.RespDuration = v.RespDuration

			}
		}
		muteRes.Unlock()
		if err := SaveResultsToFileCsv2("results.csv", new_res); err != nil {
			fmt.Println("Ошибка записи:", err)
		}
	default:
		fmt.Printf("Unknown type: %T\n", item)
	}

}

func DeleteItemService(i int) []*service.Service {

	//var srvid *service.Service
	muteSrv.RLock()

	for q, v := range new_srv {
		if v.Id == i {
			new_srv = append(new_srv[:q], new_srv[q+1:]...)
		}
	}

	muteSrv.RUnlock()

	muteSrv.Lock()
	if err := SaveServicesToFileCsv2("services.csv", new_srv); err != nil {
		fmt.Println("Ошибка записи:", err)
	}
	muteSrv.Unlock()

	return new_srv

}

func DeleteItemResult(i int) []*check.Result {

	//var srvid *service.Service
	muteRes.RLock()

	for q, v := range new_res {
		if v.Id == i {
			new_res = append(new_res[:q], new_res[q+1:]...)
		}
	}

	muteRes.RUnlock()

	muteRes.Lock()
	if err := SaveResultsToFileCsv2("results.csv", new_res); err != nil {
		fmt.Println("Ошибка записи:", err)
	}
	muteRes.Unlock()

	return new_res

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

func SaveServicesToFileCsv(path string, s *service.Service) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла")
	}
	defer file.Close()
	writeCsv := csv.NewWriter(file)
	defer writeCsv.Flush()

	str := []string{}
	str = append(str, strconv.Itoa(s.Id), s.Name, s.Url, strconv.Itoa(s.Interval), time.Now().UTC().String())
	if err := writeCsv.Write(str); err != nil {
		fmt.Println("Ошибка записи:", err)
	}

	return nil
}

func SaveServicesToFileCsv2(path string, s []*service.Service) error {

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка открытия файла")
	}
	file.Truncate(0)
	// file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// fmt.Println("Открыли файл")
	// if err != nil {
	// 	fmt.Println("Ошибка открытия файла")
	// }
	defer file.Close()
	writeCsv := csv.NewWriter(file)
	//writeCsv.Comma = '\n'
	//writeCsv.UseCRLF = true
	defer writeCsv.Flush()

	//str := []*service.Service{}
	str := []string{}

	for _, v := range s {
		//str = append(str, strconv.Itoa(v.Id), v.Name, v.Url, strconv.Itoa(v.Interval), v.Created.String())
		str = append(str, v.UnformString())
		//str = append(str, "\n")
		// fmt.Println("Записали в слайс для записи в csv")
		if err := writeCsv.Write(str); err != nil {
			fmt.Println("Ошибка записи:", err)
		}
		//writeCsv.Flush()
		str = []string{}
	}
	//str = append(str, s)
	//str = append(str, strconv.Itoa(s.Id), s.Name, s.Url, strconv.Itoa(s.Interval), time.Now().UTC().String())

	// fmt.Println("Записали в слайс для записи в csv")
	// if err := writeCsv.Write(str); err != nil {
	// 	fmt.Println("Ошибка записи:", err)
	// }

	return nil
}

func SaveResultsToFileCsv2(path string, s []*check.Result) error {

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка открытия файла")
	}
	file.Truncate(0)
	defer file.Close()
	writeCsv := csv.NewWriter(file)
	defer writeCsv.Flush()

	str := []string{}

	for _, v := range s {
		str = append(str, v.UnformString())
		if err := writeCsv.Write(str); err != nil {
			fmt.Println("Ошибка записи:", err)
		}
		//writeCsv.Flush()
		str = []string{}
	}

	return nil
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

func SaveResultsToFileCsv(path string) {
	muteRes.RLock()
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Ошибка открытия файла")
	}
	writeCsv := csv.NewWriter(file)
	defer writeCsv.Flush()

	for _, r := range new_res {
		row := []string{strconv.Itoa(r.Id),
			strconv.Itoa(r.ServiceId),
			strconv.Itoa(r.ResponseCode),
			strconv.Itoa(r.RespDuration),
			r.TimeChecked.String(),
		}
		if err := writeCsv.Write(row); err != nil {
			fmt.Println("Ошибка записи:", err)
		}
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

func LoadCsvServices(path string) error {
	//muteSrv.RLock()
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	//reader.FieldsPerRecord = 5
	records, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	//fmt.Printf("Найдено %d записей в CSV\n", cap(records))
	//fmt.Println(records)

	for i, record := range records {
		splt := strings.Split(record[0], ",")
		id, err := strconv.Atoi(splt[0])
		if err != nil {
			fmt.Printf("Ошибка парсинга ID в строке %d: %v\n", i, err)
			continue
		}

		interval, err := strconv.Atoi(splt[3])
		if err != nil {
			fmt.Printf("Ошибка парсинга Interval в строке %d: %v\n", i, err)
			continue
		}

		created, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", splt[4])
		if err != nil {
			fmt.Printf("Ошибка парсинга Created в строке %d: %v\n", i, err)
			continue
		}

		service := service.Service{
			Id:       id,
			Name:     splt[1],
			Url:      splt[2],
			Interval: interval,
			Created:  created,
		}
		new_srv = append(new_srv, &service)

	}
	//muteSrv.RUnlock()
	//fmt.Printf("Загружено %d сервисов в new_srv\n", len(new_srv))
	//fmt.Println("Послее загрузки new_srv равно", new_srv)
	return nil

}

func LoadCsvResults(path string) error {
	//muteRes.RLock()
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	for i, record := range records {
		splt := strings.Split(record[0], ",")
		id, err := strconv.Atoi(splt[0])
		if err != nil {
			fmt.Printf("Ошибка парсинга ID в строке %d: %v\n", i, err)
			continue
		}

		srvid, err := strconv.Atoi(splt[1])
		if err != nil {
			fmt.Printf("Ошибка парсинга SrvID в строке %d: %v\n", i, err)
			continue
		}

		respcode, err := strconv.Atoi(splt[2])
		if err != nil {
			fmt.Printf("Ошибка парсинга ResponseCode в строке %d: %v\n", i, err)
			continue
		}

		timecheck, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", splt[3])
		if err != nil {
			fmt.Printf("Ошибка парсинга TimeChecked в строке %d: %v\n", i, err)
			continue
		}

		duration, err := strconv.Atoi(splt[4])
		if err != nil {
			fmt.Printf("Ошибка парсинга ResponseDuration в строке %d: %v\n", i, err)
			continue
		}

		result := check.Result{
			Id:           id,
			ServiceId:    srvid,
			ResponseCode: respcode,
			TimeChecked:  timecheck,
			RespDuration: duration,
		}
		//muteRes.RUnlock()

		new_res = append(new_res, &result)

	}
	return nil

}

func SearchItemService(i int) *service.Service {

	var srvid *service.Service
	muteSrv.RLock()

	for _, v := range new_srv {
		if v.Id == i {
			srvid = v
		}
	}

	muteSrv.RUnlock()
	return srvid

}

func SearchItemResult(i int) *check.Result {

	var resid *check.Result
	muteSrv.RLock()

	for _, v := range new_res {
		if v.Id == i {
			resid = v
		}
	}

	muteSrv.RUnlock()
	return resid

}

// func ChangeItem(i int, s *service.Service) {

// 	for i, v := range new_srv {
// 		if v.Id == s.Id {

// 			sl := new_srv[i]
// 			sl.Name = s.Name
// 			sl.Url = s.Url
// 			sl.Interval = s.Interval

// 		}
// 	}

// }

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
