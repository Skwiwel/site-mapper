package mapping

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
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
		getResponse = getHTTPResponse
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

func getHTTPResponse(address, header string) (io.Reader, int, error) {
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
		network.SetExtraHTTPHeaders(prepareHeaders(header)),
		chromedp.Navigate(address),
		chromedp.OuterHTML("html", &fetchedBodyString),
	)

	// for concurrency purposes cancelling explicitly before return
	cancel()

	return strings.NewReader(fetchedBodyString), statusCode, err
}

func getHTTPResponseFast(address, header string) (io.Reader, int, error) {
	client := &http.Client{}

	request, err := createRequest(address, header)
	if err != nil {
		return nil, 0, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	return response.Body, response.StatusCode, nil
}

func createRequest(address, header string) (*http.Request, error) {
	request, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", header)
	return request, nil
}

func prepareHeaders(header string) network.Headers {
	return map[string]interface{}{
		"User-Agent": header,
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
