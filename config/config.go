package config

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
)

var (
	conf *yaml.File
)

func init() {
	LoadConf()
}

func LoadConf() {
	var err error
	conf, err = yaml.ReadFile("config.yaml")
	if err != nil {
		log.Panic(err)
	}
}

func GetItem(i string) string {
	item, _ := conf.Get(i)
	return item
}
