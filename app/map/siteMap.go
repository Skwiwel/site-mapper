package mapper

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

// SiteMap represents the hyperlink map of a website.
type SiteMap interface {
}

type siteMap struct {
	masterURL   string
	searchDepth int
	URLs        map[string]*pageNode
	urlMux      sync.Mutex
}

type pageNode struct {
	childNodes []string
	statusCode int
}

// MakeSiteMap is the siteMap constructor.
// Returns mapped site.
func MakeSiteMap(url string, searchDepth int) SiteMap {
	sm := siteMap{
		masterURL:   url,
		searchDepth: searchDepth,
		URLs:        make(map[string]*pageNode),
	}

	sm.processURL(url, searchDepth)

	return &sm
}

func (sm *siteMap) processURL(url string, depth int) {
	// check if url is mapped already
	sm.urlMux.Lock()
	if _, exists := sm.URLs[url]; exists {
		sm.urlMux.Unlock()
		return
	}
	sm.URLs[url] = &pageNode{}
	sm.urlMux.Unlock()

	// the url is not mapped yet
	// Scrape the site looking for links.
	sm.scrape(url)

}

func (sm *siteMap) scrape(url string) {
	// Create HTTP client
	client := &http.Client{}

	// Create a HTTP request before sending
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("could not create request with the url %v: %v\n", url, err)
	}
	request.Header.Set("User-Agent", "site-mapper  v0.1  no github yet :(")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Printf("could not make request to %v: %v\n", url, err)
		return
	}
	defer response.Body.Close()

	sm.URLs[url].statusCode = response.StatusCode

	document, err := html.Parse(response.Body)
	if err != nil {
		log.Printf("could not parse reponse body from %s: %v", url, err)
	}

	// Search recursively through the html document's elements to find <a href=...>'s
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attribute := range n.Attr {
				if attribute.Key == "href" {
					sm.URLs[url].childNodes = append(sm.URLs[url].childNodes, attribute.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	findLinks(document)

	for _, a := range sm.URLs[url].childNodes {
		fmt.Println(a)
	}
}
