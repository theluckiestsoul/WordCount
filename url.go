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
	taskPerRoutine := 10000
	var wg sync.WaitGroup
	splitArray := strings.Split(u.content, " ")
	remainder := len(splitArray) % taskPerRoutine
	gorroutines := len(splitArray) / taskPerRoutine
	if remainder > 0 {
		gorroutines++
	}
	
	wg.Add(gorroutines)

	count := func(s string) {
		defer wg.Done()
		u.count += strings.Count(strings.ToLower(s), u.match)
	}

	for len(splitArray) > 0 {
		if taskPerRoutine > len(splitArray) {
			taskPerRoutine = len(splitArray)
		}

		strArray := splitArray[0:taskPerRoutine]
		splitArray = splitArray[taskPerRoutine:]
		go count(strings.Join(strArray, " "))
	}
	go func() {
		wg.Wait()
		chInfo <- u
	}()
}
