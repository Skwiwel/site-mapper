package mapping

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type urlScraper struct {
	scrapedURL  url.URL
	page        *page
	pageContent *html.Node
}

func (scrapedPage *page) scrapeURLforLinks(urlScraped url.URL) error {
	scraper := urlScraper{
		scrapedURL: urlScraped,
		page:       scrapedPage,
	}

	if err := scraper.fetchPageContent(); err != nil {
		return err
	}

	scraper.findLinks()

	return nil
}

func (scraper *urlScraper) fetchPageContent() error {
	response, err := getHTTPResponse(scraper.scrapedURL.String())
	if err != nil {
		return fmt.Errorf("could not make request to %s: %w", scraper.scrapedURL.String(), err)
	}
	defer response.Body.Close()

	scraper.page.statusCode = response.StatusCode

	scraper.pageContent, err = html.Parse(response.Body)
	if err != nil {
		return fmt.Errorf("could not parse reponse body from %s: %v", scraper.scrapedURL.String(), err)
	}

	return nil
}

func getHTTPResponse(address string) (*http.Response, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create a HTTP request before sending
	request, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "site-mapper  v0.4  no github yet :(")

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
				scraper.storeChildURL(url)
			}
			break
		}
	}
}

func (scraper *urlScraper) storeChildURL(childURL *url.URL) {
	parentURL := scraper.scrapedURL
	absoluteURL := parentURL.ResolveReference(childURL)
	scraper.page.childPages[*absoluteURL] = struct{}{}
}
