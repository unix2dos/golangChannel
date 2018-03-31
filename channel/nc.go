package main

import (
	"io"
	"net"

	"fmt"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		copy(os.Stdout, conn)
		fmt.Println("done")
		done <- struct{}{}
	}()
	copy(conn, os.Stdin)
	<-done
}

func copy(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
}
