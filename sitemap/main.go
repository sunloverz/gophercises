package main

import (
	"fmt"
	"gophercises/link"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func crawl(url string, depth int, visited map[string]bool) {
	if depth <= 0 {
		return
	}
	urls := get(url)
	for _, u := range urls {
		if visited[u] {
			continue
		}
		visited[u] = true
		crawl(u, depth-1, visited)
	}
	return
}

func get(urlStr string) []string {
	resp, _ := http.Get(urlStr)
	defer resp.Body.Close()
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, value := range links {
		switch {
		case strings.HasPrefix(value.Href, "/"):
			ret = append(ret, base+value.Href)
		case strings.HasPrefix(value.Href, "http"):
			ret = append(ret, value.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

func main() {
	visited := make(map[string]bool)
	crawl("https://www.calhoun.io", 4, visited)
	links := make([]string, 0, len(visited))
	for k := range visited {
		links = append(links, k)
	}
	fmt.Println(links)
}
