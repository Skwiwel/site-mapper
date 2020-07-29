package mapping

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type urlScraper struct {
	masterURL   string
	page        *page
	pageContent *html.Node
}

func scrapeForData(urlScraped string, scrapedPage *page) {
	scraper := urlScraper{masterURL: urlScraped, page: scrapedPage}

	scraper.fetchPageContent()

	scraper.findLinks()
}

func (scraper *urlScraper) fetchPageContent() {
	response, err := getHTTPResponse(scraper.masterURL)
	if err != nil {
		log.Fatalf("could not make request to %v: %v\n", scraper.masterURL, err)
	}
	defer response.Body.Close()

	scraper.page.statusCode = response.StatusCode

	scraper.pageContent, err = html.Parse(response.Body)
	if err != nil {
		log.Printf("could not parse reponse body from %s: %v", scraper.masterURL, err)
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
	request.Header.Set("User-Agent", "site-mapper  v0.2  no github yet :(")

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
			scraper.appendLink(attribute.Val)
			break
		}
	}
}

func (scraper *urlScraper) appendLink(link string) {
	absoluteURL, err := makeURLAbsolute(link, scraper.masterURL)
	if err != nil {
		log.Printf("could not parse link %s from the document under %s: %v\n", link, scraper.masterURL, err)
		return
	}

	scraper.page.childPages = append(scraper.page.childPages, absoluteURL)
}

func makeURLAbsolute(relativeURL, masterURL string) (string, error) {
	parsedURL, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}
	if parsedURL.IsAbs() {
		return relativeURL, nil
	}
	return masterURL + relativeURL, nil
}
