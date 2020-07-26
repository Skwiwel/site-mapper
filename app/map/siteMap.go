package mapper

// SiteMap represents the hyperlink map of a website.
type SiteMap interface {
}

type siteMap struct {
	masterURL   string
	searchDepth int
	URLs        map[string]pageNode
}

type pageNode struct {
	url         string
	parentNodes map[string]pageNode
	childNodes  map[string]pageNode
}

// MakeSiteMap is the siteMap constructor.
// Returns mapped site.
func MakeSiteMap(url *string, searchDepth int) SiteMap {
	sm := siteMap{
		masterURL:   *url,
		searchDepth: searchDepth,
		URLs:        make(map[string]pageNode),
	}
	sm.URLs[sm.masterURL] = makePageNode(&sm.masterURL)
	return &sm
}

func makePageNode(url *string) pageNode {
	return pageNode{
		url:         *url,
		parentNodes: make(map[string]pageNode),
		childNodes:  make(map[string]pageNode),
	}
}
