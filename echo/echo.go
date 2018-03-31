package main

import (
	"bufio"
	"io"
	"net"
	"strings"
	"time"
)

func main() {

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		str := scanner.Text()

		go func(str string) {
			io.WriteString(conn, "\t"+strings.ToUpper(str)+"\n")
			time.Sleep(time.Second)
			io.WriteString(conn, "\t"+str+"\n")
			time.Sleep(time.Second)
			io.WriteString(conn, "\t"+strings.ToLower(str)+"\n")
		}(str)
	}
}
