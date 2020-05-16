package config

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	config := Cfg

	fmt.Printf("%#v\n", *config)
	fmt.Printf("%#v\n", *config.Email)
	fmt.Printf("%#v\n", *config.Cookie)
}
