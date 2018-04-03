package main

import (
	"io"
	"net"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8001")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go MustCopy(conn, os.Stdin)
	MustCopy(os.Stdout, conn)
}

func MustCopy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
}
