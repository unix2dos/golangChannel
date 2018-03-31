package main

import "fmt"

func main() {
	naturnals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 0; x < 10; x++ {
			naturnals <- x
		}
		close(naturnals)
	}()

	go func() {
		for v := range naturnals {
			squares <- v * v
		}
		close(squares)
	}()

	for v := range squares {
		fmt.Println(v)
	}
}
