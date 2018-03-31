package main

import (
	"fmt"

	"os"

	"golang.org/x/net/html"
)

func crawl(url string) []string {
	fmt.Println(url)

	list, err := extract(url)
	if err != nil {

	}
	return list
}

func extract(url string) (list []string, err error) {

	return
}

func main() {

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		panic(err)
	}

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}

}
