package models

import (
	"testing"

	"github.com/jdxj/bilibili/config"
)

func TestAPIResponse_LoginInfo(t *testing.T) {
	gu := config.GetUsers()
	b, err := NewBiliBili("test", gu[0].Cookie)
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = b.tryCheckSign()
	if err != nil {
		t.Fatalf("%s", err)
	}
}
