/**
 * main.go
 *
 * @author : Cyw
 * @email  : rose2099.c@gmail.com
 * @created: 2017/6/12 下午6:32
 * @logs   :
 *
 */
package main

import (
	"flag"
	"go-sign/libary"
	_ "go-sign/site"
)

func main() {
	flag.Parse()

	for _, site := range flag.Args() {
		libary.Open(site)
	}
}