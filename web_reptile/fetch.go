package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	for _, v := range os.Args[1:] {
		resp, err := http.Get(v)
		if err != nil {
			continue
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
	}

}
