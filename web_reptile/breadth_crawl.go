package main

import (
	"fmt"
	"net/http"
	"os"

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

func Crawl(url string) (list []string) {
	fmt.Println(url)
	list, err := Extract(url)
	if err != nil {
		return
	}
	return
}

func BreadthFirst(f func(url string) []string, worklist []string) {

	exist := make(map[string]bool)

	for len(worklist) > 0 {
		bak := worklist
		worklist = nil

		for _, url := range bak {

			if !exist[url] {

				exist[url] = true

				worklist = append(worklist, f(url)...)
			}
		}
	}

}

func main() {

	BreadthFirst(Crawl, os.Args[1:])
}
