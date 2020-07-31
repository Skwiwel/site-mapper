package mapping

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	tu "github.com/skwiwel/site-mapper/test/testutil"
	"golang.org/x/net/html"
)

func Test_urlScraper_storeChildURL(t *testing.T) {
	t.Parallel()
	type args struct {
		childURL *url.URL
	}
	tests := []struct {
		name                string
		scraper             *urlScraper
		args                args
		presentInChildPages url.URL
		wantInChildPages    url.URL
	}{
		{
			name:                "relative url",
			scraper:             &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			args:                args{tu.URLParseSkipError("http://foo.com", t)},
			presentInChildPages: url.URL{},
			wantInChildPages:    *tu.URLParseSkipError("http://foo.com", t),
		},
		{
			name:                "absolute url",
			scraper:             &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			args:                args{tu.URLParseSkipError("/users", t)},
			presentInChildPages: url.URL{},
			wantInChildPages:    *tu.URLParseSkipError("http://example.com/users", t),
		},
		{
			name:                "repeated child",
			scraper:             &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			args:                args{tu.URLParseSkipError("/users", t)},
			presentInChildPages: *tu.URLParseSkipError("http://example.com/users", t),
			wantInChildPages:    *tu.URLParseSkipError("http://example.com/users", t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.scraper.page.childPages[tt.presentInChildPages] = struct{}{}
			tt.scraper.storeChildURL(tt.args.childURL)
			if _, found := tt.scraper.page.childPages[tt.wantInChildPages]; !found {
				t.Errorf("storeChildURL() did not properly store relative url")
				return
			}
		})
	}
}

func Test_urlScraper_parseNodeContent(t *testing.T) {
	type args struct {
		n *html.Node
	}
	tests := []struct {
		name             string
		scraper          *urlScraper
		htmlString       string
		wantInChildPages url.URL
		wantEmpty        bool
	}{
		{
			name:             "is link",
			scraper:          &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			htmlString:       `<a href="http://example.com">TargetMarker</a>`,
			wantInChildPages: *tu.URLParseSkipError("http://example.com", t),
			wantEmpty:        false,
		},
		{
			name:             "no link",
			scraper:          &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			htmlString:       `<div>TargetMarker</div>`,
			wantInChildPages: *tu.URLParseSkipError("http://example.com", t),
			wantEmpty:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fake html
			str := tt.htmlString
			r := strings.NewReader(str)
			pageContent, _ := html.Parse(r)

			// document traverse
			var findLinks func(*html.Node)
			findLinks = func(n *html.Node) {
				if n.FirstChild != nil && n.FirstChild.Data == "TargetMarker" {
					tt.scraper.parseNodeContent(n)
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					findLinks(c)
				}
			}
			findLinks(pageContent)

			// check result
			found := false
			empty := true
			for url := range tt.scraper.page.childPages {
				empty = false
				if reflect.DeepEqual(url, tt.wantInChildPages) {
					found = true
				}
			}
			if tt.wantEmpty != empty {
				if empty {
					t.Errorf("parseNodeContent() haven't found and stored the link from %s", tt.htmlString)
				} else {
					t.Errorf("parseNodeContent() incorrectly saved an url from %s", tt.htmlString)
				}
				return
			}
			if !empty && !found {
				t.Errorf("parseNodeContent() saved an incorrect url from %s", tt.htmlString)
			}
		})
	}
}

func Test_urlScraper_findLinks(t *testing.T) {
	tests := []struct {
		name           string
		scraper        *urlScraper
		htmlString     string
		wantChildPages map[url.URL]struct{}
	}{
		{
			name:       "has link",
			scraper:    &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			htmlString: `<a href="http://example.com"></a>`,
			wantChildPages: map[url.URL]struct{}{
				*tu.URLParseSkipError("http://example.com", t): struct{}{},
			},
		},
		{
			name:           "no link",
			scraper:        &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			htmlString:     `<div></div>`,
			wantChildPages: map[url.URL]struct{}{},
		},
		{
			name:       "multiple links",
			scraper:    &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			htmlString: `<a href="http://example1.com"></a><a href="http://example2.com"></a><a href="http://example3.com"></a>`,
			wantChildPages: map[url.URL]struct{}{
				*tu.URLParseSkipError("http://example1.com", t): struct{}{},
				*tu.URLParseSkipError("http://example2.com", t): struct{}{},
				*tu.URLParseSkipError("http://example3.com", t): struct{}{},
			},
		},
		{
			name:       "relative links",
			scraper:    &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			htmlString: `<a href="http://example1.com"></a><a href="/users"></a><a href="/shop/target"></a>`,
			wantChildPages: map[url.URL]struct{}{
				*tu.URLParseSkipError("http://example1.com", t):            struct{}{},
				*tu.URLParseSkipError("http://example.com/users", t):       struct{}{},
				*tu.URLParseSkipError("http://example.com/shop/target", t): struct{}{},
			},
		},
		{
			name:       "repeating links",
			scraper:    &urlScraper{*tu.URLParseSkipError("http://example.com", t), &page{childPages: make(map[url.URL]struct{})}, &html.Node{}},
			htmlString: `<a href="http://example1.com"></a><a href="/users"></a><a href="/users"></a>`,
			wantChildPages: map[url.URL]struct{}{
				*tu.URLParseSkipError("http://example1.com", t):      struct{}{},
				*tu.URLParseSkipError("http://example.com/users", t): struct{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fake html
			str := tt.htmlString
			r := strings.NewReader(str)
			tt.scraper.pageContent, _ = html.Parse(r)

			tt.scraper.findLinks()

			// check result
			if !reflect.DeepEqual(tt.scraper.page.childPages, tt.wantChildPages) {
				t.Errorf("findLinks() created an incorrect child link map\n   got: %#v\nwanted: %#v",
					tt.scraper.page.childPages,
					tt.wantChildPages)
			}
		})
	}
}

func Test_urlScraper_scrapeURLforLinks(t *testing.T) {
	tests := []struct {
		name       string
		wantError  bool
		wantStatus int
		wantURLs   map[string]struct{} // URLs should be relative, where appropriate
		htmlString string
	}{
		{
			name:       "normal page",
			wantStatus: http.StatusOK,
			wantURLs: map[string]struct{}{
				"http://example1.com": struct{}{},
				"/users":              struct{}{},
			},
			htmlString: `<a href="http://example1.com"></a><a href="/users"></a><a href="/users"></a>`,
		},
		{
			name:       "empty page",
			wantStatus: http.StatusOK,
			wantURLs:   map[string]struct{}{},
			htmlString: `<body></body>`,
		},
		{
			name:       "no server",
			wantError:  true,
			wantURLs:   map[string]struct{}{},
			htmlString: `<body></body>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetURL := *tu.URLParseSkipError("http://addressToNowhere", t)
			if tt.wantStatus == http.StatusOK {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=UTF-8")
					fmt.Fprintln(w, tt.htmlString)
				}))
				defer mockServer.Close()

				targetURL = *tu.URLParseSkipError(mockServer.URL, t)
			}

			resultPage := &page{childPages: make(map[url.URL]struct{})}
			err := resultPage.scrapeURLforLinks(targetURL)
			if (err != nil) != tt.wantError {
				if err != nil {
					t.Errorf("scrapeURLforLinks() improperly returned an error: %v", err)
				} else {
					t.Errorf("scrapeURLforLinks() did not return an error when it should.")
				}
			}

			// fix the wantURLs map to include only absolute URL strings
			for urlRaw := range tt.wantURLs {
				if url := tu.URLParseSkipError(urlRaw, t); !url.IsAbs() {
					delete(tt.wantURLs, urlRaw)
					absoluteURL := targetURL.String() + urlRaw
					tt.wantURLs[absoluteURL] = struct{}{}
				}
			}

			if resultPage.statusCode != tt.wantStatus {
				t.Errorf("scrapeURLforLinks() saved a wrong status code %v; should be %v", resultPage.statusCode, tt.wantStatus)
			}
			for url := range resultPage.childPages {
				if _, found := tt.wantURLs[url.String()]; !found {
					t.Errorf("scrapeURLforLinks() found too many URLs - additional URL: %s\n", url.String())
				}
			}
			for urlRaw := range tt.wantURLs {
				if _, found := resultPage.childPages[*tu.URLParseSkipError(urlRaw, t)]; !found {
					t.Errorf("scrapeURLforLinks() failed to save all URLs - not found: %s\n", urlRaw)
				}
			}
		})
	}
}
