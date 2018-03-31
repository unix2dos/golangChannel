package tcp

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	copy(os.Stdout, conn)
}

func copy(dst io.Writer, src io.Reader) {

	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}
