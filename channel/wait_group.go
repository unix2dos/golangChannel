package main

import (
	"os"
	"sync"
)

func makeThumb1(filenames <-chan string) (total int64) {

	size := make(chan int64)
	var wg sync.WaitGroup

	for v := range filenames {
		wg.Add(1)

		go func(v string) {
			defer wg.Done()

			thumb := v //假设是个函数

			file, err := os.Stat(thumb)
			if err != nil {
				return
			}
			size <- file.Size()

		}(v)
	}

	//必须基于size循环的并发
	go func() {
		wg.Wait()
		close(size)
	}()

	//1. 如果wait在main goroutine 循环前,  size的东西取不走, 阻塞死

	for v := range size {
		total += v
	}

	//2. 如果wait在main goroutine 循环后,  循环永远不会终止, 没人关闭size

	return
}
