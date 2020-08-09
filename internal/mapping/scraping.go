package mapping

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
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
	responseString, statusCode, err := getHTTPResponse(scraper.scrapedURL.String())
	if err != nil {
		return fmt.Errorf("could not make request to %s: %w", scraper.scrapedURL.String(), err)
	}

	scraper.page.statusCode = statusCode

	scraper.pageContent, err = html.Parse(strings.NewReader(responseString))
	if err != nil {
		return fmt.Errorf("could not parse reponse body from %s: %v", scraper.scrapedURL.String(), err)
	}

	return nil
}

func getHTTPResponse(address string) (string, int, error) {
	chromeContext, cancel := chromedp.NewContext(context.Background())

	var fetchedBodyString string
	var statusCode int

	chromedp.ListenTarget(chromeContext, func(event interface{}) {
		switch responseReceivedEvent := event.(type) {
		case *network.EventResponseReceived:
			response := responseReceivedEvent.Response
			if response.URL == address || response.URL == address+"/" {
				statusCode = int(response.Status)
			}
		}
	})

	err := chromedp.Run(chromeContext,
		network.Enable(),
		network.SetExtraHTTPHeaders(prepareHeaders()),
		chromedp.Navigate(address),
		chromedp.OuterHTML("html", &fetchedBodyString),
	)

	// for concurrency purposes cancelling explicitly before return
	cancel()

	return fetchedBodyString, statusCode, err
}

func prepareHeaders() network.Headers {
	return map[string]interface{}{
		"User-Agent": "site-mapper  v0.5, github: https://gtihub.com/skwiwel/site-mapper",
	}
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
