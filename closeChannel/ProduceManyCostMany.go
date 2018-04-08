package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumProduces = 100
	const NumCosters = 10

	dataCh := make(chan int, 100) //数据
	stopCh := make(chan struct{})
	stopByNameCh := make(chan string, 1) //1是为了裁判goroutine没准备好就收到通知了
	var stopName string

	//裁判goroutine
	go func() {
		stopName = <-stopByNameCh
		close(stopCh)
	}()

	wg := sync.WaitGroup{}
	wg.Add(NumCosters)

	//produce
	for i := 0; i < NumProduces; i++ {
		go func(i int) {
			for {
				value := rand.Intn(MaxRandomNumber)
				if value == 0 {
					select {
					case stopByNameCh <- strconv.Itoa(i): //加select是为了不阻塞
					default:
					}
					return
				}

				select { //这边写是为了尽早
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(i)
	}

	//cost
	for i := 0; i < NumCosters; i++ {
		go func(i int) {
			defer wg.Done()

			for {

				select { //这边写是为了尽早
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == MaxRandomNumber-1 {
						select {
						case stopByNameCh <- strconv.Itoa(i): //加select是为了不阻塞
						default:
						}
						return
					}
					log.Println(value)
				}
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("---> " + stopName)
}
