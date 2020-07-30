package util

import (
	"log"
	"net/url"
)

// VerifyAndParseURL checks if the address is in a good format and prepends http:// if needed
func VerifyAndParseURL(address string) url.URL {
	url, err := url.Parse(address)
	if err != nil {
		log.Fatalf("could not parse address %s: %v", address, err)
	}
	urlE := urlExtended{url}
	urlE.prependHTTPScheme()
	return *urlE.URL
}

type urlExtended struct {
	*url.URL
}

func (url *urlExtended) prependHTTPScheme() {
	if !url.IsAbs() {
		url.Scheme = "http"
	}
}
