package main

import (
	"flag"
	"log"
	"net/url"

	mapper "github.com/skwiwel/site-mapper/app/map"
)

var urlString = flag.String("url", "", "The site to map")

func main() {
	flag.Parse()
	if *urlString == "" {
		log.Fatal("Error: The --url flag must be set.")
	}

	url, err := url.Parse(*urlString)
	if err != nil {
		log.Fatal(err)
	}

	// Add http:// if needed
	if !url.IsAbs() && url.Scheme == "" {
		url.Scheme = "http"
	}

	siteMap := mapper.MakeSiteMap(url.String(), 1)
	siteMap.Print()
}
