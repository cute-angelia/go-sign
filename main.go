package main

import (
	"flag"
	"go-sign/libary"
	_ "go-sign/site"
	"log"
)

func main() {
	flag.Parse()

	log.Print(libary.Drivers())

	for _, site := range flag.Args() {
		libary.Open(site)
	}

	conf := libary.GetItem("signs")
	log.Print(conf)
}
