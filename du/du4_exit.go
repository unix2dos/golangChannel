package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var progress = flag.Bool("v", false, "show progress")

var limit = make(chan struct{}, 20)

var done = make(chan struct{})

func isExist() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func print(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1fGB\n", nfiles, float64(nbytes)/1e9)
}

func Dir(name string, size chan int64, n *sync.WaitGroup) (err error) {
	defer n.Done()
	limit <- struct{}{}
	defer func() { <-limit }()
	if isExist() {
		return
	}

	fileinfo, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	for _, v := range fileinfo {

		if v.IsDir() {
			n.Add(1)
			go Dir(filepath.Join(name, v.Name()), size, n)

		} else {

			size <- v.Size()
		}
	}

	return
}

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
	var n sync.WaitGroup
	for _, v := range args {
		n.Add(1)
		go Dir(v, channle_size, &n)
	}

	go func() {
		n.Wait()
		close(channle_size)
	}()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	var nfiles, nbytes int64
loop:
	for {
		select {
		case <-done:
			for range channle_size {
			}
			//return //排空, return
			panic("11") //小技巧, 看其他goroutine是否正确
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
