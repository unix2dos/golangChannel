package main

import (
	"net/http"
	"os"

	"fmt"

	"golang.org/x/net/html"
)

var depth int

func ParseStart(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		depth++
	}

}

func ParseEnd(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}

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

func main() {
	for _, v := range os.Args[1:] {
		resp, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		doc, err := html.Parse(resp.Body)
		if err != nil {
			continue
		}
		ForeachNode(doc, ParseStart, ParseEnd)
	}
}
