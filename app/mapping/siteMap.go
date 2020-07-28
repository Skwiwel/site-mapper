package mapping

import (
	"sync"
)

// MappedLocation represents the hyperlink map of a website.
type MappedLocation interface {
	Print()
}

type siteMap struct {
	masterURL string
	depth     int
	URLs      sync.Map // string -> *page
}

type page struct {
	childPages []string
	statusCode int
}

// MapSite is the siteMap constructor.
// Returns mapped site.
func MapSite(address string, depth int) MappedLocation {
	sm := siteMap{
		masterURL: address,
		depth:     depth,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	processURL(address, depth, &sm.URLs, &wg)
	wg.Wait()

	return &sm
}

func processURL(address string, depth int, URLs *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	// Check if the address is already mapped
	processedPage := &page{}
	_, loaded := URLs.LoadOrStore(address, processedPage)
	if loaded { // the page was already processed before
		return
	}

	// the url is not mapped yet
	// Scrape the site looking for links.
	scrapeForData(address, processedPage)
	// Process child URLs
	for _, c := range processedPage.childPages {
		depth--
		if depth < 0 {
			return
		}
		wg.Add(1)
		go processURL(c, depth, URLs, wg)
	}

}
