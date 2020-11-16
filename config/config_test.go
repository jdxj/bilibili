package config

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	e := GetEmail()
	fmt.Printf("%#v\n", *e)
	us := GetUsers()
	for _, v := range us {
		fmt.Printf("%#v\n", *v)
	}
}

func TestStructCopy(t *testing.T) {
}
