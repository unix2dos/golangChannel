package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 3)
		cancel()
	}()

	fmt.Println(A(ctx))

	time.Sleep(time.Hour)
}

func A(ctx context.Context) string {

	go fmt.Println(B(ctx))

	for {
		select {
		case <-ctx.Done():
			return "A Done"
		}
	}
	return ""
}

func B(ctx context.Context) string {

	go fmt.Println(C(ctx))
	for {
		select {
		case <-ctx.Done():
			return "B Done"
		}
	}
	return ""
}

func C(ctx context.Context) string {
	for {
		select {
		case <-ctx.Done():
			return "C Done"
		}
	}
	return ""
}
