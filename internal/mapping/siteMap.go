package mapping

import (
	"fmt"
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
func MapSite(url url.URL, depth int) (MappedLocation, error) {
	sm := siteMap{
		masterURL: url,
		depth:     depth,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	if err := sm.processURL(sm.masterURL, depth, &wg); err != nil {
		return nil, fmt.Errorf("could not map the site under %s: %w", url.String(), err)
	}
	wg.Wait()

	return &sm, nil
}

func (sm *siteMap) processURL(url url.URL, depth int, wg *sync.WaitGroup) error {
	defer wg.Done()

	// Check if the address is already mapped
	processedPage := &page{}
	_, alreadyProcessed := sm.URLs.LoadOrStore(url.String(), processedPage)
	if alreadyProcessed {
		return nil
	}

	if err := processedPage.scrapeURLforLinks(url); err != nil {
		return fmt.Errorf("could not scrape url %s: %w", url.String(), err)
	}

	sm.processChildURLs(processedPage, depth, wg)

	return nil
}

func (sm *siteMap) processChildURLs(parentPage *page, depth int, wg *sync.WaitGroup) {
	for childURL := range parentPage.childPages {
		depth--
		if depth < 0 {
			return
		}
		wg.Add(1)
		// ignoring errors from child pages
		go sm.processURL(childURL, depth, wg)
	}
	return
}
