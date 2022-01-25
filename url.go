package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type urlInfo struct {
	url     string
	match   string
	count   int
	content string
}

func (u *urlInfo) readContent() {
	res, err := http.Get(u.url)
	if err != nil {
		log.Fatal(err)
	}
	html, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	u.content = string(html)
}

func (u urlInfo) String() string {
	return fmt.Sprintf("%v - %v \n", u.url, u.count)
}

func (u urlInfo) countMatch(chInfo chan<- urlInfo) {
	var wg sync.WaitGroup
	gorroutines := len(u.content) / 10000
	wg.Add(gorroutines)
	runes := []rune(u.content)

	count := func(s string) {
		defer wg.Done()
		u.count += strings.Count(strings.ToLower(s), u.match)
	}

	for i := 0; i < gorroutines; i++ {
		startIndex, endIndex := 10000*i, 10000*(i+1)
		temp := runes[startIndex:endIndex]
		go count(string(temp))
	}
	go func() {
		wg.Wait()
		chInfo <- u
	}()
}
