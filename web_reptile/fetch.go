package main

import (
	"fmt"
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
		resp.Body.Close()
		if err != nil {
			continue
		}
		fmt.Println(string(bytes))
	}

}
