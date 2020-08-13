package learningtests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

func Test_chromedp_Run_Simple(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		htmlString string
	}{
		{
			name: "normal page",
			htmlString: `<html>
			  <head></head>
			  <body>
				  <a href="/3">3</a>
				  <a href="/1">1</a>
			  </body>
			</html>`,
		},
		{
			name:       "empty page",
			htmlString: `<html><head></head><body></body></html>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=UTF-8")
				fmt.Fprintln(w, tt.htmlString)
			}))
			defer mockServer.Close()

			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			var fetchedBodyString string
			err := chromedp.Run(ctx,
				chromedp.Navigate(mockServer.URL),
				chromedp.OuterHTML("html", &fetchedBodyString),
			)
			if err != nil {
				t.Errorf("chromedp.Run() failed to fetch the site content: %v", err)
			}

			// The fetched document has an added newline before </body>.
			// It's generally harmless. Here it's removed for the sake of the test comparison.
			fetchedBodyString = trimNewlineBeforeBodyClose(fetchedBodyString)

			//log.Println(tt.htmlString)
			//log.Println(fetchedBodyString)

			testDocument, err := html.Parse(strings.NewReader(tt.htmlString))
			if err != nil {
				t.Errorf("html.Parse() failed to parse the test htmlString: %v", err)
			}
			fetchedDocument, err := html.Parse(strings.NewReader(fetchedBodyString))
			if err != nil {
				t.Errorf("html.Parse() failed to parse the fetched html body: %v", err)
			}

			var testDocumentBody *html.Node
			var fetchedDocumentBody *html.Node

			var crawlForBody func(*html.Node, **html.Node)
			crawlForBody = func(node *html.Node, targetNode **html.Node) {
				if node.Type == html.ElementNode && node.Data == "body" {
					*targetNode = node
					return
				}
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					crawlForBody(child, targetNode)
				}
			}
			crawlForBody(testDocument, &testDocumentBody)
			crawlForBody(fetchedDocument, &fetchedDocumentBody)

			if !reflect.DeepEqual(testDocumentBody, fetchedDocumentBody) {
				t.Errorf("chromedp.Run() fetched site content in an unsupported format\n   got: %#v\nwanted: %#v",
					fetchedDocumentBody,
					testDocumentBody)
			}
		})
	}
}

func trimNewlineBeforeBodyClose(str string) string {
	pos := strings.Index(str, "\n</body>")
	if pos != -1 {
		str = str[0:pos] + str[pos+1:]
	}
	return str
}

func Test_chromedp_Run_Extended(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		htmlString string
	}{
		{
			name: "normal page",
			htmlString: `<html>
			  <head></head>
			  <body>
				  <a href="/3">3</a>
				  <a href="/1">1</a>
			  </body>
			</html>`,
		},
		{
			name:       "empty page",
			htmlString: `<html><head></head><body></body></html>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=UTF-8")
				fmt.Fprintln(w, tt.htmlString)
			}))
			defer mockServer.Close()

			chromeContext, cancelContext := chromedp.NewContext(context.Background())
			defer cancelContext()

			var fetchedBodyString string
			var statusCode int64
			var statusCodeMux sync.Mutex

			chromedp.ListenTarget(chromeContext, func(event interface{}) {
				switch responseReceivedEvent := event.(type) {
				case *network.EventResponseReceived:
					response := responseReceivedEvent.Response
					if response.URL == mockServer.URL+"/" {
						statusCodeMux.Lock()
						statusCode = response.Status
						statusCodeMux.Unlock()
					}
				}
			})

			runError := chromedp.Run(
				chromeContext,
				network.Enable(),
				chromedp.Navigate(mockServer.URL),
				chromedp.OuterHTML("html", &fetchedBodyString),
			)

			if runError != nil {
				t.Errorf("chromedp.Run() failed to fetch the site content: %v", runError)
			}

			statusCodeMux.Lock()
			defer statusCodeMux.Unlock()
			if statusCode != http.StatusOK {
				t.Errorf("chromedp.Run() fetched with status: %v", statusCode)
			}

			// The fetched document has an added newline before </body>.
			// It's generally harmless. Here it's removed for the sake of the test comparison.
			fetchedBodyString = trimNewlineBeforeBodyClose(fetchedBodyString)

			//log.Println(tt.htmlString)
			//log.Println(fetchedBodyString)

			testDocument, err := html.Parse(strings.NewReader(tt.htmlString))
			if err != nil {
				t.Errorf("html.Parse() failed to parse the test htmlString: %v", err)
			}
			fetchedDocument, err := html.Parse(strings.NewReader(fetchedBodyString))
			if err != nil {
				t.Errorf("html.Parse() failed to parse the fetched html body: %v", err)
			}

			var testDocumentBody *html.Node
			var fetchedDocumentBody *html.Node

			var crawlForBody func(*html.Node, **html.Node)
			crawlForBody = func(node *html.Node, targetNode **html.Node) {
				if node.Type == html.ElementNode && node.Data == "body" {
					*targetNode = node
					return
				}
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					crawlForBody(child, targetNode)
				}
			}
			crawlForBody(testDocument, &testDocumentBody)
			crawlForBody(fetchedDocument, &fetchedDocumentBody)

			if !reflect.DeepEqual(testDocumentBody, fetchedDocumentBody) {
				t.Errorf("chromedp.Run() fetched site content in an unsupported format\n   got: %#v\nwanted: %#v",
					fetchedDocumentBody,
					testDocumentBody)
			}
		})
	}
}
