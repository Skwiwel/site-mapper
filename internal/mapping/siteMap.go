package mapping

import (
	"net/url"
	"sync"
)

// MappedLocation represents the hyperlink map of a website.
type MappedLocation interface {
	Print()
}

type siteMap struct {
	masterURL url.URL
	depth     int
	URLs      sync.Map // string -> *page
}

type page struct {
	childPages map[url.URL]struct{}
	statusCode int
}

// MapSite is the siteMap constructor.
// Returns mapped site.
func MapSite(url url.URL, depth int) MappedLocation {
	sm := siteMap{
		masterURL: url,
		depth:     depth,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	sm.processURL(sm.masterURL, depth, &wg)
	wg.Wait()

	return &sm
}

func (sm *siteMap) processURL(url url.URL, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Check if the address is already mapped
	processedPage := &page{}
	_, alreadyProcessed := sm.URLs.LoadOrStore(url.String(), processedPage)
	if alreadyProcessed {
		return
	}

	processedPage.scrapeURLforLinks(url)
	// Process child URLs
	sm.processChildURLs(processedPage, depth, wg)

}

func (sm *siteMap) processChildURLs(parentPage *page, depth int, wg *sync.WaitGroup) {
	for childURL := range parentPage.childPages {
		depth--
		if depth < 0 {
			return
		}
		wg.Add(1)
		go sm.processURL(childURL, depth, wg)
	}
}
