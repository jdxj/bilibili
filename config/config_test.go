package config

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	config := readConfig("../config.json")
	fmt.Printf("%#v\n", *config.Email)
	for _, v := range config.Cookies {
		fmt.Printf("%#v\n", *v)
	}
}

func TestStructCopy(t *testing.T) {
	c1 := &Cookie{}
	c2 := *c1

	fmt.Printf("%p\n", c1)
	fmt.Printf("%p\n", &c2)
}
