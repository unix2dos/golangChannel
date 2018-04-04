package main

import (
	"image"
	"sync"
)

var icons map[string]image.Image

func loadIcons() {
	icons = map[string]image.Image{
		"spades.png":   new(image.Image),
		"hearts.png":   new(image.Image),
		"diamonds.png": new(image.Image),
		"clubs.png":    new(image.Image),
	}
}

func Icon(name string) image.Image {
	if icons == nil {
		loadIcons()
	}
	return icons[name]
}

//互斥
var mux sync.Mutex

func Icon_Safe(name string) image.Image { //会影响并发访问
	mux.Lock()
	defer mux.Unlock()
	if icons == nil {
		loadIcons()
	}
	return icons[name]
}

//读并发访问
var mux2 sync.RWMutex

func Icon_RSafe(name string) image.Image {

	mux2.RLock()

	if icons != nil {
		icon := icons[name]
		mux2.RUnlock()
		return icon
	}

	mux2.RUnlock()
	mux2.Lock()

	if icons != nil { //此处为什么一定要check呢,因为上面2个锁之间有间隙, 可能又被其他goroutine load了
		loadIcons()
	}

	icon := icons[name]
	mux2.Unlock()
	return icon
}

//once
var loadOnce sync.Once

func Icon_Once(name string) image.Image {
	loadOnce.Do(loadIcons)
	return icons[name]
}

func main() {

}
