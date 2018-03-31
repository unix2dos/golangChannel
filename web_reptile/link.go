package main

import "golang.org/x/net/html"

//获取一个link的所有子link, 可能会无限递归
func visit(links []string, n *html.Node) []string {

	if n.Type == html.ElementNode && n.Data == "a" {
		//又是个链接
		for _, v := range n.Attr {
			if v.Key == "href" {
				links = append(links, v.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}
