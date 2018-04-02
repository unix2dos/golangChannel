package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

func Dir(name string, size chan int64) (err error) {

	fileinfo, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	for _, v := range fileinfo {

		if v.IsDir() {
			Dir(filepath.Join(name, v.Name()), size)

		} else {

			size <- v.Size()
		}
	}

	return
}

var progress = flag.Bool("v", false, "show progress")

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var tick <-chan time.Time
	if *progress {
		tick = time.Tick(500 * time.Millisecond)
	}

	channle_size := make(chan int64)
	go func() {

		for _, v := range args {
			Dir(v, channle_size)
		}
		close(channle_size)

	}()

	var nfiles, nbytes int64

loop:
	for {
		select {
		case size, ok := <-channle_size:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += size
		case <-tick:
			print(nfiles, nbytes)
		}
	}

	print(nfiles, nbytes)

}
func print(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1fGB\n", nfiles, float64(nbytes)/1e9)
}
