package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/skwiwel/site-mapper/app/mapping"
)

var address = flag.String("url", "", "The site to map")

func main() {
	flag.Parse()
	if *address == "" {
		log.Fatal("Error: The --url flag must be set.")
	}

	url, err := url.Parse(*address)
	if err != nil {
		log.Fatal(err)
	}

	// Add http:// if needed
	if !url.IsAbs() && url.Scheme == "" {
		url.Scheme = "http"
	}

	siteMap := mapping.MapSite(url.String(), 1)
	siteMap.Print()
}
