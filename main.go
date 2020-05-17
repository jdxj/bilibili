package main

import (
	"github.com/jdxj/bilibili/client"
)

func main() {
	c := client.NewClient()
	c.Start()
}
