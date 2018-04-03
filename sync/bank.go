package main

import (
	"fmt"
	"time"
)

var chan_desposits = make(chan int)
var chan_balances = make(chan int)

func Deposit(amount int) {
	chan_desposits <- amount
}

func Balance() int {
	return <-chan_balances
}

func teller() {
	var sum int
	for {
		select {
		case amount := <-chan_desposits:
			sum += amount
		case chan_balances <- sum:
		}
	}
}

func init() {
	go teller()
}

func main() {

	go func() {
		Deposit(10)
	}()

	go func() {
		Deposit(5)
	}()

	fmt.Println(Balance())

	time.Sleep(time.Second)
}
