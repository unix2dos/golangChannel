package main

import "fmt"

func main() {

loop:
	for i := 0; i < 100; i++ {
		if i == 3 {
			break loop
		}
		fmt.Println(i)
	}
	fmt.Println("finish")
}
