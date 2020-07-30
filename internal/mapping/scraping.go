package mapping

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type urlScraper struct {
	scrapedURL  url.URL
	page        *page
	pageContent *html.Node
}

func (scrapedPage *page) scrapeURLforLinks(urlScraped url.URL) {
	scraper := urlScraper{scrapedURL: urlScraped, page: scrapedPage}
	scraper.fetchPageContent()

	scraper.findLinks()
}

func (scraper *urlScraper) fetchPageContent() {
	response, err := getHTTPResponse(scraper.scrapedURL.String())
	if err != nil {
		log.Fatalf("could not make request to %s: %v\n", scraper.scrapedURL.String(), err)
	}
	defer response.Body.Close()

	scraper.page.statusCode = response.StatusCode

	scraper.pageContent, err = html.Parse(response.Body)
	if err != nil {
		log.Printf("could not parse reponse body from %s: %v", scraper.scrapedURL.String(), err)
	}
}

func getHTTPResponse(address string) (*http.Response, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create a HTTP request before sending
	request, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "site-mapper  v0.3  no github yet :(")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (scraper *urlScraper) findLinks() {
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		scraper.parseNodeContent(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	findLinks(scraper.pageContent)
}

func (scraper *urlScraper) parseNodeContent(n *html.Node) {
	if n.Type != html.ElementNode || n.Data != "a" { // "a", as in <a href=...>
		return
	}
	for _, attribute := range n.Attr {
		if attribute.Key == "href" {
			if url, err := url.Parse(attribute.Val); err == nil {
				scraper.appendChildURL(url)
			} else {
				log.Printf("could not parse link %s found on %s\n", attribute.Val, scraper.scrapedURL.String())
			}
			break
		}
	}
}

func (scraper *urlScraper) appendChildURL(childURL *url.URL) {
	parentURL := scraper.scrapedURL
	absoluteURL := parentURL.ResolveReference(childURL)
	scraper.page.childPages = append(scraper.page.childPages, *absoluteURL)
}
