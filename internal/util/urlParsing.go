package util

import (
	"net/url"
)

// VerifyAndParseURL checks if the address is in a good format and prepends http:// if needed
func VerifyAndParseURL(address string) (url.URL, error) {
	parsedURL, err := url.Parse(address)
	if err != nil {
		return url.URL{}, err
	}
	urlE := urlExtended{parsedURL}
	urlE.prependHTTPScheme()
	return *urlE.URL, nil
}

type urlExtended struct {
	*url.URL
}

func (url *urlExtended) prependHTTPScheme() {
	if !url.IsAbs() {
		url.Scheme = "http"
		newURL, _ := url.Parse(url.String()) // by reparsing the url.URL is cleaned
		*url = urlExtended{newURL}
	}
}
