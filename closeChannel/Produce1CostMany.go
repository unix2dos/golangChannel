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
	const NumCosters = 100

	dataCh := make(chan int, 100) //数据

	wg := sync.WaitGroup{}
	wg.Add(NumCosters)

	//produce
	go func() {
		for {
			if value := rand.Intn(MaxRandomNumber); value == 0 {
				// the only produce can close the channel safely.
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	//cost
	for i := 0; i < NumCosters; i++ {
		go func() {
			defer wg.Done()

			for value := range dataCh {
				log.Println(value)
			}

		}()
	}

	wg.Wait()
}
