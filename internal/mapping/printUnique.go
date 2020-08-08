package mapping

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/skwiwel/site-mapper/internal/colored"
	"github.com/skwiwel/site-mapper/internal/status"
)

// PrintUnique prints the unique URLs found in the search
func (sm *siteMap) PrintUnique() {
	tempString := sm.masterURL.String() + "\n"

	tempString += printSubAddresses(&sm.URLs, &sm.masterURL)

	// print the depth of the search
	tempString += fmt.Sprintf("\nmapped depth: %v\n", sm.depth)

	fmt.Print(tempString)
}

type mapRange interface {
	Range(f func(key, value interface{}) bool)
}

func printSubAddresses(URLs mapRange, masterURL fmt.Stringer) string {
	tempString := ""
	URLs.Range(func(key interface{}, value interface{}) bool {
		address := key.(string)
		page := value.(*page)

		if address == masterURL.String() {
			return true
		}

		tempString += "   "

		address, _ = getURLWithoutScheme(address)

		switch page.statusCode {
		case http.StatusOK:
			tempString += address
		case status.Unprocessed:
			tempString += addressWarning(address, page.queries, "connection failed")
		case status.OutOfScope:
			tempString += addressNote(address, page.queries, "out of scope")
		default:
			tempString += addressWarning(address, page.queries, page.statusCode)
		}
		tempString += "\n"
		return true
	})
	return tempString
}

func getURLWithoutScheme(address string) (string, error) {
	parsedURL, err := url.Parse(address)
	if err != nil {
		return address, fmt.Errorf("could not parse url string: %w", err)
	}
	return fmt.Sprintf("%s%s", parsedURL.Hostname(), parsedURL.RequestURI()), nil
}

func addressFormatting(addressFormat, commentFormat func(...interface{}) string) func(interface{}, bool, interface{}) string {
	return func(address interface{}, queriesPresent bool, comment interface{}) string {
		queriesMarker := ""
		if queriesPresent {
			queriesMarker = "?..."
		}
		return addressFormat(fmt.Sprint(address)) +
			colored.Grey(queriesMarker) +
			" " +
			commentFormat(fmt.Sprint(comment))
	}
}

var (
	addressWarning = addressFormatting(colored.Grey, colored.Yellow)
	addressNote    = addressFormatting(colored.White, colored.Grey)
)
