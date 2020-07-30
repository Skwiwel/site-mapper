package util

import (
	"net/url"
	"reflect"
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
		name        string
		urlModified *urlExtended
		urlWanted   *urlExtended
	}{
		{"no scheme", &urlExtended{urlParseSkipError("example.com", t)}, &urlExtended{urlParseSkipError("http://example.com", t)}},
		{"has http scheme", &urlExtended{urlParseSkipError("http://example.com", t)}, &urlExtended{urlParseSkipError("http://example.com", t)}},
		{"has https scheme", &urlExtended{urlParseSkipError("https://example.com", t)}, &urlExtended{urlParseSkipError("https://example.com", t)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.urlModified.prependHTTPScheme()
			if !reflect.DeepEqual(tt.urlModified, tt.urlWanted) {
				t.Errorf("VerifyAndParseURL()\n got: %#v\nwant: %#v", tt.urlModified, tt.urlWanted)
			}
		})
	}
}

func TestVerifyAndParseURL(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    url.URL
		wantErr bool
	}{
		{"correct and absolute url", args{"http://example.com"}, *urlParseSkipError("http://example.com", t), false},
		{"correct but without scheme", args{"example.com"}, *urlParseSkipError("http://example.com", t), false},
		{"incorrect - only port", args{":80"}, url.URL{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyAndParseURL(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyAndParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyAndParseURL()\n got: %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
