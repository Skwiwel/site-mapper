package main

import (
	"flag"
	"log"

	"github.com/skwiwel/site-mapper/internal/mapping"
	"github.com/skwiwel/site-mapper/internal/util"
)

func main() {
	address := flag.String("url", "", "The site to map")
	depth := flag.Int("depth", 2, "depth of the search")
	fast := flag.Bool("fast", false, `the mapper has 2 modes of operation:
	Normally pages are fetched using a headless chrome process. This ensures that pages are fully loaded with javascript parsed before processing.
	When using --fast the pages are fetched without parsing any javascript. This can lead to an incomplete mapping of a site. The fast takes only a fraction of the time taken by normal operation.`)
	flag.Parse()
	if *address == "" {
		log.Fatal("Error: The --url flag must be set.")
	}

	url, err := util.VerifyAndParseURL(*address)
	if err != nil {
		log.Fatalf("could not parse the given url %s: %v", *address, err)
	}

	siteMap, err := mapping.MapSite(url, *depth, *fast)
	if err != nil {
		log.Fatal(err)
	}

	siteMap.PrintUnique()
}
