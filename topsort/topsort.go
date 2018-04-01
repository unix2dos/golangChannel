package main

import (
	"fmt"
	"sort"
)

var prereqs = map[string][]string{
	"algorithms":            {"data structures"},
	"calculus":              {"linear algebra"},
	"compilers":             {"data structures", "formal languages", "computer organization"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func Sort() (list []string) {

	var keys []string
	for key := range prereqs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	exist := make(map[string]bool)

	var visit func(item string)
	visit = func(item string) {

		if !exist[item] {

			exist[item] = true

			for _, v := range prereqs[item] {
				visit(v)
			}

			list = append(list, item)
		}
	}

	for _, v := range keys {
		visit(v)
	}

	//exsit := make(map[string]bool)
	//var visitAll func(items []string)
	//visitAll = func(items []string) {
	//
	//	for _, v := range items {
	//
	//		if exsit[v] == false {
	//			exsit[v] = true
	//			visitAll(prereqs[v])
	//			list = append(list, v)
	//		}
	//
	//	}
	//}
	//visitAll(keys)

	return
}

func main() {
	for _, v := range Sort() {
		fmt.Println(v)
	}
}
