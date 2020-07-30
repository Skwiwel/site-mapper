package util

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/skwiwel/site-mapper/test/testutil"
)

func Test_urlExtended_prependHTTPScheme(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		urlModified *urlExtended
		urlWanted   *urlExtended
	}{
		{"no scheme", &urlExtended{testutil.URLParseSkipError("example.com", t)}, &urlExtended{testutil.URLParseSkipError("http://example.com", t)}},
		{"has http scheme", &urlExtended{testutil.URLParseSkipError("http://example.com", t)}, &urlExtended{testutil.URLParseSkipError("http://example.com", t)}},
		{"has https scheme", &urlExtended{testutil.URLParseSkipError("https://example.com", t)}, &urlExtended{testutil.URLParseSkipError("https://example.com", t)}},
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
	t.Parallel()
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    url.URL
		wantErr bool
	}{
		{"correct and absolute url", args{"http://example.com"}, *testutil.URLParseSkipError("http://example.com", t), false},
		{"correct but without scheme", args{"example.com"}, *testutil.URLParseSkipError("http://example.com", t), false},
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
