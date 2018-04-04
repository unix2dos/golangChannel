package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"sync"

	"github.com/prometheus/common/log"
)

func incomingURLs() []string {
	return []string{
		"https://www.baidu.com",
		"http://www.163.com",
		"http://www.weibo.com",
		"http://www.qq.com",
		"https://www.baidu.com",
		"http://www.163.com",
		"http://www.weibo.com",
		"http://www.qq.com",
	}
}

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type Func func(string) (interface{}, error)

type result struct {
	val interface{}
	err error
}

type entry struct {
	ready chan struct{}
	res   result
}

type request struct {
	key      string
	response chan<- result
}

type Memo struct {
	requests chan request
}

func NewMemo(f Func) *Memo {
	memo := &Memo{make(chan request)}
	go memo.server(f)
	return memo
}

func (m *Memo) Close() {
	close(m.requests)
}

func (m *Memo) server(f Func) {
	cache := make(map[string]*entry)

	for req := range m.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			e.res.val, e.res.err = f(req.key)
			close(e.ready)
		} else {
			<-e.ready
		}
		req.response <- e.res
	}
}

func (m *Memo) Get(key string) (interface{}, error) {

	response := make(chan result)
	m.requests <- request{key, response}
	res := <-response
	return res.val, res.err

}

func main() {

	m := NewMemo(httpGetBody)

	var group sync.WaitGroup
	for _, url := range incomingURLs() {

		group.Add(1)
		go func(url string) {
			defer group.Done()

			start := time.Now()
			val, err := m.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("url = %s  time = %s  bytes = %d\n", url, time.Since(start), len(val.([]byte)))
		}(url)
	}

	group.Wait()
}
