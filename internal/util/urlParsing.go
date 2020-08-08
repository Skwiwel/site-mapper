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

// CleanAndCheckForQuery removes url's # references and query contents.
// Returns cleaned address and a boolean indicating wether query was present.
func CleanAndCheckForQuery(address url.URL) (url.URL, bool, error) {
	address.Fragment = ""
	query := false
	if address.RawQuery != "" {
		query = true
		address.RawQuery = ""
	}
	cleanedURL, err := url.Parse(address.String())
	if err != nil { // badly constructed url.URL
		return *cleanedURL, false, err
	}
	return *cleanedURL, query, nil
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
