package examples

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

type Fetcher interface {
	// Fetch returns the body of URL and a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	defer wg.Done()
	if depth <= 0 {
		return
	}

	// Don't fetch the same URL twice.
	if _, loaded := monitor.LoadOrStore(url, struct{}{}); loaded {
		// fmt.Printf("... Skipping %s\n", url)
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		// Fetch URLs in parallel.
		go Crawl(u, depth-1, fetcher)
	}
}

var wg sync.WaitGroup

// var monitor = NewMonitor() // faster and efficient (1 B/op, 0 allocs/op) than sync.Map
var monitor sync.Map // slower and less efficient (17 B/op, 1 allocs/op) than our custom implementation using RWMutex

func TestWebCrawler(t *testing.T) {
	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher)
	wg.Wait()
	fmt.Println("Finished crawling the webpages")
}

func BenchmarkWebCrawler(b *testing.B) {
	before := runtime.NumGoroutine()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			go Crawl("https://golang.org/", 4, fetcher)
		}
	})
	wg.Wait()
	fmt.Printf("Finished crawling the webpages. Go routines Before: %d and After: %d\n", before, runtime.NumGoroutine())
}

type Monitor struct {
	mux  sync.RWMutex
	urls map[string]struct{}
}

func NewMonitor() *Monitor {
	return &Monitor{
		urls: make(map[string]struct{}),
	}
}

func (m *Monitor) LoadOrStore(url string, value interface{}) (alwaysNil interface{}, loaded bool) {
	m.mux.RLock()
	if _, ok := m.urls[url]; ok {
		m.mux.RUnlock()
		return nil, true
	}
	m.mux.RUnlock()
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.urls[url]; ok {
		return nil, true
	}
	m.urls[url] = struct{}{}
	return nil, false
}

type fakeResult struct {
	body string
	urls []string
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
