package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/html"
)

var fetched map[string]bool

// Crawl uses findLinks to recursively crawl
// pages starting with url, to a maximum of depth.
// Non-concurrent version.
//func Crawl(url string, depth int) {
//	// TODO: Fetch URLs in parallel.
//
//	if depth < 0 {
//		return
//	}
//	urls, err := findLinks(url)
//	if err != nil {
//		// fmt.Println(err)
//		return
//	}
//	fmt.Printf("found: %s\n", url)
//	fetched[url] = true
//	for _, u := range urls {
//		if !fetched[u] {
//			Crawl(u, depth-1)
//		}
//	}
//	return
//}

// result holds data to be passed between goroutines
type result struct {
	url   string
	urls  []string
	err   error
	depth int
}

// Crawl uses findLinks to recursively crawl
// pages starting with url, to a maximum of depth.
// Concurrent version.
func Crawl(url string, depth int) {
	// channel to send the result struct on
	results := make(chan *result)

	// fetch function will call findLinks and send results to the results channel.
	fetch := func(url string, depth int) {
		urls, err := findLinks(url)
		results <- &result{url, urls, err, depth}
	}

	// start a goroutine to call fetch
	go fetch(url, depth)

	// record that we've fetched the url passed to this function
	// so we don't get it again
	fetched[url] = true

	for fetching := 1; fetching > 0; fetching-- {
		// get the results off the channel
		res := <-results
		if res.err != nil {
			continue
		}

		fmt.Printf("found %s\n", res.url)

		// Keep calling fetch until depth is 0.
		// Depth is decremented each time fetch is called.
		if res.depth > 0 {
			for _, u := range res.urls {
				// Only proceed if we haven't already seen this url
				if !fetched[u] {
					fetching++
					go fetch(u, res.depth-1)
					fetched[u] = true
				}
			}
		}
	}

	close(results)
}

func main() {
	// Set the max number of CPUs that can be executing
	// to the number of CPUs on the machine.
	// Redundant to set this as GOMAXPROCS defaults to max
	// number of CPUs anyeay.
	runtime.GOMAXPROCS(runtime.NumCPU())
	fetched = make(map[string]bool)
	now := time.Now()
	Crawl("http://andcloud.io", 2)
	fmt.Println("time taken:", time.Since(now))
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

// visit appends to links each link found in n, and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
