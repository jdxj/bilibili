package main

import (
	"github.com/jdxj/bilibili/client"
	"github.com/jdxj/bilibili/config"
)

func main() {
	c := client.NewClient(config.Cfg.Cookie)
	c.Start()
}
