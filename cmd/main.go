package main

import (
	"flag"
	"log"

	"github.com/skwiwel/site-mapper/internal/mapping"
	"github.com/skwiwel/site-mapper/internal/util"
)

func main() {
	address := flag.String("url", "", "The site to map")
	flag.Parse()
	if *address == "" {
		log.Fatal("Error: The --url flag must be set.")
	}

	url := util.VerifyAndParseURL(*address)

	siteMap := mapping.MapSite(url, 1)
	siteMap.Print()
}
