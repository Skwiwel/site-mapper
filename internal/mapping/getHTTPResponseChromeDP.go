package mapping

import (
	"context"
	"io"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func getHTTPResponseChromeDP(address, header string) (io.Reader, int, error) {
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

func prepareHeaders(header string) network.Headers {
	return map[string]interface{}{
		"User-Agent": header,
	}
}
