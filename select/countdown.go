package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for i := 3; i > 0; i-- {

		select {
		case <-tick.C:
			fmt.Println(i)
		case <-abort:
			fmt.Println("abort")
			return
		}

	}

	fmt.Println("launch")
}
