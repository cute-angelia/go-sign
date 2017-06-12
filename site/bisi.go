package site

import (
	"go-sign/libary"
	"log"
)

type BisiDriver struct{}

func (d BisiDriver) Run() (rest libary.Rest, e error) {

	log.Print("goog")
	rest.Code = 1
	return rest, e
}

func init()  {
	libary.Register("bisi", BisiDriver{})
}