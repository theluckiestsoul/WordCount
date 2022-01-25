package main

//avinash@go-mmt.com
import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var sg sync.WaitGroup
	chInfo := make(chan urlInfo)

	urls := []urlInfo{{url: "https://makemytrip.com", match: "makemytrip"}, {url: "https://goibibo.com", match: "goibibo"}}

	sg.Add(len(urls))

	for _, url := range urls {
		url.readContent()
		go url.countMatch(chInfo)
	}

	go func() {
		defer close(chInfo)
		for v := range chInfo {
			fmt.Println(v)
			sg.Done()
		}
	}()

	sg.Wait()
	printMemUsage()

}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
