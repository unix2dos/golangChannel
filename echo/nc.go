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
	defer conn.Close()

	go copy(os.Stdout, conn)
	copy(conn, os.Stdin)
}

func copy(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
}
