package main

import (
	"io"
	"net"

	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}

	go copy(conn, os.Stdin)
	copy(os.Stdout, conn)
}

func copy(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
}
