package mapping

import (
	"fmt"
	"io"
	"net/url"

	"golang.org/x/net/html"
)

var header string = "site-mapper  v0.5, github: https://gtihub.com/skwiwel/site-mapper"

type urlScraper struct {
	scrapedURL  url.URL
	page        *page
	pageContent *html.Node
}

var getResponse func(address, header string) (io.Reader, int, error)

func (scrapedPage *page) scrapeURLforLinks(urlScraped url.URL, fast bool) error {
	scraper := urlScraper{
		scrapedURL: urlScraped,
		page:       scrapedPage,
	}

	setResponseGetter(fast)

	if err := scraper.fetchPageContent(); err != nil {
		return err
	}

	scraper.findLinks()

	return nil
}

func setResponseGetter(fast bool) {
	if !fast {
		getResponse = getHTTPResponseChromeDP
	} else {
		getResponse = getHTTPResponseFast
	}
}

func (scraper *urlScraper) fetchPageContent() error {
	responseString, statusCode, err := getResponse(scraper.scrapedURL.String(), header)
	if err != nil {
		return fmt.Errorf("could not make request to %s: %w", scraper.scrapedURL.String(), err)
	}

	scraper.page.statusCode = statusCode

	scraper.pageContent, err = html.Parse(responseString)
	if err != nil {
		return fmt.Errorf("could not parse reponse body from %s: %v", scraper.scrapedURL.String(), err)
	}

	return nil
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
