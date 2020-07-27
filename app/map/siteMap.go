package mapper

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"
)

// SiteMap represents the hyperlink map of a website.
type SiteMap interface {
	Print()
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

	var wg sync.WaitGroup
	wg.Add(1)
	sm.processURL(url, searchDepth, &wg)
	wg.Wait()

	return &sm
}

func (sm *siteMap) processURL(url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()
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
	// Process child URLs
	for _, c := range sm.URLs[url].childNodes {
		depth--
		if depth < 0 {
			return
		}
		wg.Add(1)
		go sm.processURL(c, depth, wg)
	}

}

func (sm *siteMap) scrape(urlScraped string) {
	// Create HTTP client
	client := &http.Client{}

	// Create a HTTP request before sending
	request, err := http.NewRequest("GET", urlScraped, nil)
	if err != nil {
		log.Fatalf("could not create request with the url %v: %v\n", urlScraped, err)
	}
	request.Header.Set("User-Agent", "site-mapper  v0.2  no github yet :(")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Printf("could not make request to %v: %v\n", urlScraped, err)
		return
	}
	defer response.Body.Close()

	sm.URLs[urlScraped].statusCode = response.StatusCode

	document, err := html.Parse(response.Body)
	if err != nil {
		log.Printf("could not parse reponse body from %s: %v", urlScraped, err)
	}

	// Search recursively through the html document's elements to find <a href=...>'s
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attribute := range n.Attr {
				if attribute.Key == "href" {
					urlFound := attribute.Val
					parsedURL, err := url.Parse(urlFound)
					if err != nil {
						log.Printf("could not parse url %s: %v\n", urlFound, err)
						continue
					}
					if !parsedURL.IsAbs() {
						urlFound = urlScraped + urlFound
					}
					sm.URLs[urlScraped].childNodes = append(sm.URLs[urlScraped].childNodes, urlFound)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	findLinks(document)
}

func (sm *siteMap) Print() {
	var color func(string) func(...interface{}) string
	color = func(colorString string) func(...interface{}) string {
		return func(args ...interface{}) string {
			return fmt.Sprintf(colorString,
				fmt.Sprint(args...))
		}
	}
	warning := color("\033[33;5m%s\033[0m")
	grey := color("\033[33;90m%s\033[0m")

	var tempString string

	tempString += "\n"
	tempString += sm.masterURL + "\n"
	// Print each unique url found.
	// Accesses information potentially actively modified by some goroutines.
	// Mutex isn't locked so as to prioritize their work over this Print() func.
	for urlString, node := range sm.URLs {
		if urlString == sm.masterURL {
			continue
		}
		// Clear http://
		parsedURL, err := url.Parse(urlString)
		if err != nil {
			log.Fatal("could not parse url string: ", err)
		}
		urlString = fmt.Sprintf("   %s%s", parsedURL.Hostname(), parsedURL.RequestURI())
		if node.statusCode == http.StatusOK {
			tempString += urlString
		} else {
			tempString += grey(urlString)
			tempString += " "
			tempString += fmt.Sprintf(warning(node.statusCode))
		}
		tempString += "\n"
	}

	// print the depth of the search
	tempString += fmt.Sprintf("\nmapped depth: %v\n", sm.searchDepth)

	fmt.Println(tempString)
}
