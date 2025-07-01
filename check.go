package main

import "fmt"

func main() {

	var x int
	var y int

	fmt.Print("Введите количество столбцов: ")
	fmt.Scan(&x)
	fmt.Print("Введите количество строк: ")
	fmt.Scan(&y)

	j := 0
	for j < y {

		if j%2 == 0 {
			i := 0
			for i < x {
				if i%2 == 0 {
					fmt.Print(string('#'))
				} else {
					fmt.Print(string(' '))
				}
				i++
			}
		} else {
			i := 0
			for i < x {
				if i%2 == 0 {
					fmt.Print(string(' '))
				} else {
					fmt.Print(string('#'))
				}
				i++
			}
		}
		fmt.Println()

		j++
	}

}
