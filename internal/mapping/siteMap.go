package mapping

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/skwiwel/site-mapper/internal/status"
	"github.com/skwiwel/site-mapper/internal/util"
)

// MappedLocation represents the hyperlink map of a website.
type MappedLocation interface {
	PrintUnique()
}

type siteMap struct {
	masterURL url.URL
	depth     int
	URLs      sync.Map // string -> *page
}

type page struct {
	childPages map[url.URL]struct{}
	statusCode int
	queries    bool
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

func (sm *siteMap) processURL(address url.URL, depth int, wg *sync.WaitGroup) error {
	defer wg.Done()

	address, queryPresent, err := util.CleanAndCheckForQuery(address)
	if err != nil {
		return fmt.Errorf("could not clean an url: %w", err)
	}

	processedPage, alreadyProcessed := sm.createPageForProcessing(address.String())
	if alreadyProcessed {
		return nil
	}
	// Assumes that if multiple queries per uri are present then they are always present.
	// Ex. youtube's /watch?... is never just /watch with no query
	processedPage.queries = queryPresent

	if !isInSearchScope(&address, &sm.masterURL) {
		processedPage.statusCode = status.OutOfScope
		return nil
	}

	if err := processedPage.scrapeURLforLinks(address); err != nil {
		return fmt.Errorf("could not scrape url %s: %w", address.String(), err)
	}

	sm.processChildURLs(processedPage, depth, wg)

	return nil
}

func (sm *siteMap) createPageForProcessing(address string) (*page, bool) {
	processedPage := &page{childPages: map[url.URL]struct{}{}}
	_, alreadyProcessed := sm.URLs.LoadOrStore(address, processedPage)
	if alreadyProcessed {
		return nil, true
	}
	return processedPage, false
}

func isInSearchScope(address fmt.Stringer, masterAddress fmt.Stringer) bool {
	return strings.Contains(address.String(), masterAddress.String())
}

func (sm *siteMap) processChildURLs(parentPage *page, depth int, wg *sync.WaitGroup) {
	depth--
	if depth < 0 {
		return
	}
	for childURL := range parentPage.childPages {
		wg.Add(1)
		// ignoring errors from child pages
		go sm.processURL(childURL, depth, wg)
	}
	return
}
