package main

import (
	"flag"
	"go-sign/libary"
	"log"
)



func main() {
	flag.Parse()

	for _, site := range flag.Args() {
		log.Print(site)
	}

	conf := libary.GetItem("signs")
	log.Print(conf)
}
