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

	*address = util.VerifyAndRepairURL(*address)

	siteMap := mapping.MapSite(*address, 1)
	siteMap.Print()
}
