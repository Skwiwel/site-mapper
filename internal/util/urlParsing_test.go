package util

import (
	"net/url"
	"testing"
)

func urlParseSkipError(urlRaw string, t *testing.T) *url.URL {
	url, err := url.Parse(urlRaw)
	if err != nil {
		t.Errorf("could not parse %s: %v", urlRaw, err)
	}
	return url
}

func Test_urlExtended_prependHTTPScheme(t *testing.T) {
	tests := []struct {
		name string
		url  *urlExtended
		want string
	}{
		{"no scheme", &urlExtended{urlParseSkipError("example.com", t)}, "http"},
		{"has http scheme", &urlExtended{urlParseSkipError("http://example.com", t)}, "http"},
		{"has https scheme", &urlExtended{urlParseSkipError("https://example.com", t)}, "https"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.url.prependHTTPScheme()
			if got := tt.url.Scheme; got != tt.want {
				t.Errorf("prependHTTPScheme() url has scheme %s, want %s", got, tt.want)
			}
		})
	}
}

func TestVerifyAndParseURL(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want url.URL
	}{
		{"correct and absolute url", args{"http://example.com"}, *urlParseSkipError("http://example.com", t)},
		{"correct but without scheme", args{"example.com"}, *urlParseSkipError("http://example.com", t)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyAndParseURL(tt.args.address); got.String() != tt.want.String() {
				t.Errorf("VerifyAndParseURL() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}
