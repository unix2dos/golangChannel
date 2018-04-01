package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func ForeachNode1(n *html.Node, start, end func(n *html.Node)) {
	if start != nil {
		start(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForeachNode1(c, start, end)
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

	ForeachNode1(doc, vistnode, nil)

	return
}

func main() {
	for _, v := range os.Args[1:] {

		list, err := Extract(v)
		if err != nil {
			continue
		}

		for _, vv := range list {
			fmt.Println(vv)
		}

	}
}
