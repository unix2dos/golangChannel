package main

import (
	"bufio"
	"fmt"
	"net"
)

var enter = make(chan chan string)
var exit = make(chan chan string)
var message = make(chan string)

func main() {

	listener, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleWithConn(conn)
	}
}

func broadcaster() {

	var clients = make(map[chan string]bool)

	for {
		select {
		case msg := <-message:
			for ch := range clients {
				ch <- msg
			}
		case ch := <-enter:
			clients[ch] = true
		case ch := <-exit:
			delete(clients, ch)
			close(ch)
		}
	}

}

func handleWithConn(conn net.Conn) {
	defer conn.Close()

	var ch = make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	ch <- "hello, weclome " + who
	message <- who + " has arrived"
	enter <- ch

	scan := bufio.NewScanner(conn)
	for scan.Scan() {
		message <- who + ": " + scan.Text()
	}

	exit <- ch
	message <- who + " has leaved"

}

func clientWriter(conn net.Conn, ch <-chan string) {

	for str := range ch {
		fmt.Fprintln(conn, str)
	}
}
