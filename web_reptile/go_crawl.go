package main

import (
	"fmt"
	"net/http"

	"os"

	"net/url"

	"golang.org/x/net/html"
)

func ForeachNode(n *html.Node, start, end func(n *html.Node)) {
	if start != nil {
		start(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForeachNode(c, start, end)
	}

	if end != nil {
		end(n)
	}
}

func Extract(url string) (list []string, err error) {

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	vistnode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			for _, v := range n.Attr {
				if v.Key != "href" {
					continue

				}
				link, err := resp.Request.URL.Parse(v.Val)
				if err != nil {
					continue
				}
				list = append(list, link.String())
			}
		}
	}

	ForeachNode(doc, vistnode, nil)

	return
}

var limit = make(chan struct{}, 20)

func Crawl(strUrl string) (list []string) {

	str, _ := url.QueryUnescape(strUrl)
	fmt.Println(str)

	limit <- struct{}{}
	list, err := Extract(strUrl)
	<-limit
	if err != nil {
		return
	}
	return
}

func main() {

	worklist := make(chan []string)
	go func() {
		worklist <- os.Args[1:]
	}()

	var n int
	n++

	exist := make(map[string]bool)

	for ; n > 0; n-- {

		list := <-worklist

		for _, url := range list {

			if !exist[url] {

				exist[url] = true

				n++

				go func(url string) {

					worklist <- Crawl(url)

				}(url)
			}
		}
	}
}
