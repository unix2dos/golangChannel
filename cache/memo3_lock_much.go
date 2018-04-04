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

type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]result
}

func NewMemo(fun Func) *Memo {
	return &Memo{f: fun, cache: make(map[string]result)}
}

func (m *Memo) Get(url string) (interface{}, error) {
	m.mu.Lock()
	res, ok := m.cache[url]
	m.mu.Unlock()

	if !ok {
		res.val, res.err = m.f(url)
		m.mu.Lock()
		m.cache[url] = res
		m.mu.Unlock()
	}
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
