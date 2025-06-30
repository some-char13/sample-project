package main

import "fmt"

func main() {

	var x int
	var y int

	fmt.Print("Введите количество строк: ")
	fmt.Scanln(&x)
	fmt.Print("Введите количество столбцов: ")
	fmt.Scanln(&y)

	z := y / 2

	for i := 0; i < x; i++ {
		for q := 0; q < z; q++ {
			if i%2 == 0 {
				fmt.Print(string('#') + string(' '))
			} else {
				fmt.Print(string(' ') + string('#'))

			}
		}
		fmt.Println()
	}
}
