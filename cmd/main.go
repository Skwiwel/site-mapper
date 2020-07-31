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
	flag.Parse()
	if *address == "" {
		log.Fatal("Error: The --url flag must be set.")
	}

	url, err := util.VerifyAndParseURL(*address)
	if err != nil {
		log.Fatalf("could not parse the given url %s: %v", *address, err)
	}

	siteMap, err := mapping.MapSite(url, *depth)
	if err != nil {
		log.Fatal(err)
	}

	siteMap.Print()
}
