package mapping

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func scrapeForData(urlScraped string, scrapedPage *page) {
	response, err := getHTTPResponse(urlScraped)
	if err != nil {
		log.Fatalf("could not make request to %v: %v\n", urlScraped, err)
	}
	defer response.Body.Close()

	scrapedPage.statusCode = response.StatusCode

	document, err := html.Parse(response.Body)
	if err != nil {
		log.Printf("could not parse reponse body from %s: %v", urlScraped, err)
	}

	// Search recursively through the html document's elements to find <a href=...>'s
	findLinksInHTML(document, &scrapedPage.childPages, urlScraped)

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

func findLinksInHTML(document *html.Node, addresses *[]string, masterAddress string) {
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if err := appendIfLink(n, addresses, masterAddress); err != nil {
			log.Printf("failed to process a HTML node of %s: %v", masterAddress, err)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	findLinks(document)
}

func appendIfLink(n *html.Node, addresses *[]string, masterAddress string) error {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attribute := range n.Attr {
			if attribute.Key == "href" {
				if err := appendAddress(addresses, attribute.Val, masterAddress); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func appendAddress(addresses *[]string, address, masterAddress string) error {
	// concatenate the address to masterAddress if it's relative to it
	parsedURL, err := url.Parse(address)
	if err != nil {
		return err
	}
	if !parsedURL.IsAbs() {
		address = masterAddress + address
	}

	*addresses = append(*addresses, address)
	return nil
}
