package main

import (
	"fmt"
	"net/http"
)

func main() {
	url := "https://ya.ru"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Response code:", resp.Status)
	} else if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		fmt.Println("this is not ok", resp.StatusCode)
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}
