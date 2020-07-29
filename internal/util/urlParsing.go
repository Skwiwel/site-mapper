package util

import (
	"log"
	"net/url"
)

// VerifyAndRepairURL checks if the address is in a good format and prepends http:// if needed
func VerifyAndRepairURL(address string) string {
	url, err := url.Parse(address)
	if err != nil {
		log.Fatalf("could not parse address %s: %v", address, err)
	}
	urlE := urlExtended{url}
	urlE.prependHTTPScheme()
	return urlE.String()
}

type urlExtended struct {
	*url.URL
}

func (url *urlExtended) prependHTTPScheme() {
	if url.shouldHaveScheme() {
		url.Scheme = "http"
	}
}

func (url *urlExtended) shouldHaveScheme() bool {
	return !url.IsAbs() && url.Scheme == ""
}
