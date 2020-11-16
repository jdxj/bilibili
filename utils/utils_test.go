package utils

import (
	"testing"

	"github.com/jdxj/bilibili/config"
)

func TestSendMessage(t *testing.T) {
	ge := config.GetEmail()
	err := SendMessage(ge.User, "abc", "def")
	if err != nil {
		t.Fatalf("%s", err)
	}
}
