package mapping

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

func color(colorString string) func(...interface{}) string {
	return func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
}

var (
	warning = color("\033[33;5m%s\033[0m")
	grey    = color("\033[33;90m%s\033[0m")
)

// Print prints the siteMap structure to console
func (sm *siteMap) Print() {
	tempString := sm.masterURL + "\n"

	tempString += printSubAddresses(&sm.URLs, sm.masterURL)

	// print the depth of the search
	tempString += fmt.Sprintf("\nmapped depth: %v\n", sm.depth)

	fmt.Print(tempString)
}

func printSubAddresses(URLs *sync.Map, masterURL string) string {
	tempString := ""
	URLs.Range(func(key interface{}, value interface{}) bool {
		address := key.(string)
		page := value.(*page)

		if address == masterURL {
			return true
		}

		tempString += "   "
		address = getURLWithoutScheme(address)

		if page.statusCode == http.StatusOK || page.statusCode == 0 {
			tempString += address
		} else {
			tempString += grey(address)
			tempString += " "
			tempString += fmt.Sprintf(warning(page.statusCode))
		}
		tempString += "\n"
		return true
	})
	return tempString
}

func getURLWithoutScheme(address string) string {
	parsedURL, err := url.Parse(address)
	if err != nil {
		log.Fatal("could not parse url string: ", err)
	}
	return fmt.Sprintf("%s%s", parsedURL.Hostname(), parsedURL.RequestURI())
}
