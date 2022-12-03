package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/nakajima/pingsy/pkg/config"
	"github.com/nakajima/pingsy/pkg/reporter"
)

func reportError(url string, message string) {
	fmt.Printf("Error checking URL: %s\n", message)
	reporter.ReportError(url, message)
}

func check(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(url)

	if err != nil {
		reportError(url, err.Error())
		return
	}

	if response.StatusCode != 200 {
		reportError(url, fmt.Sprintf("Status for %s was %d", url, response.StatusCode))
	}

	fmt.Println(url)
}

func main() {
	var wg sync.WaitGroup

	urls := config.URLS()

	for _, url := range urls {
		wg.Add(1)
		go check(url, &wg)
	}

	wg.Wait()
	fmt.Println("Done")
}
