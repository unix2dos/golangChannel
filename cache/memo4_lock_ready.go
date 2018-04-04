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

type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]*entry
}

func NewMemo(fun Func) *Memo {
	return &Memo{f: fun, cache: make(map[string]*entry)}
}

func (m *Memo) Get(url string) (interface{}, error) {
	m.mu.Lock()
	e := m.cache[url]
	if e == nil {
		e = &entry{ready: make(chan struct{})}
		m.cache[url] = e
		m.mu.Unlock()

		e.res.val, e.res.err = m.f(url)
		close(e.ready)
	} else {
		m.mu.Unlock()
		<-e.ready
	}
	return e.res.val, e.res.err
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
