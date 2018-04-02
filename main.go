package main

import (
	"fmt"
	"net/url"
)

func main() {
	var urlStr string = "https://unix2dos.github.io/2018/03/07/%E5%8C%BA%E5%9D%97%E9%93%BE%E4%BB%8B%E7%BB%8D/#%E5%88%86%E7%B1%BB"

	l3, _ := url.QueryUnescape(urlStr)
	fmt.Println(l3)

	a, _ := url.ParseRequestURI(urlStr)
	fmt.Println(a)
	fmt.Println(a.Query())
}
