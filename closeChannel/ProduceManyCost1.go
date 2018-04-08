package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumProduces = 100

	dataCh := make(chan int, 100) //数据
	stopCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(1)

	//produce
	for i := 0; i < NumProduces; i++ {
		go func() {
			for {
				select {
				case <-stopCh:
					return
				default:
					value := rand.Intn(MaxRandomNumber)
					dataCh <- value
				}

			}
		}()
	}

	//cost
	go func() {
		defer wg.Done()
		for value := range dataCh {
			if value == 0 {
				close(stopCh)
				return
			}

			log.Println(value)
		}

	}()

	wg.Wait()
}
