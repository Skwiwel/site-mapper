package main

import (
	"flag"
	"fmt"
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

	siteMap := mapper.MakeSiteMap(url.String(), 1)
	fmt.Println(siteMap)
}
