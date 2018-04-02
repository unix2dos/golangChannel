package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	channle_size := make(chan int64)

	go func() {

		for _, v := range args {
			Dir(v, channle_size)
		}
		close(channle_size)

	}()

	var nfiles, nbytes int64
	for size := range channle_size {
		nfiles++
		nbytes += size
	}
	print(nfiles, nbytes)

}
func print(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1fGB\n", nfiles, float64(nbytes)/1e9)
}
