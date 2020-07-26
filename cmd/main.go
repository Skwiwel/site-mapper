package main

import (
	"flag"
	"log"

	"github.com/skwiwel/site-mapper/app/rpattern"
)

var siteURL = flag.String("url", "", "The site to map")

func main() {
	flag.Parse()
	if *siteURL == "" {
		log.Fatal("Error: The --url flag must be set.")
	}
	if match := rpattern.WebURL.MatchString(*siteURL); !match {
		log.Fatal("Error: Incorrect url.")
	}
}
