package mapping

import (
	"fmt"
	"net/http"
	"net/url"
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
			break
		case 0:
			tempString += grey(address)
			tempString += " "
			tempString += fmt.Sprintf(warning("connection failed"))
			break
		default:
			tempString += grey(address)
			tempString += " "
			tempString += fmt.Sprintf(warning(page.statusCode))
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
