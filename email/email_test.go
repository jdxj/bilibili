package email

import (
	"testing"

	"github.com/jdxj/bilibili/config"
)

func TestNewEmail(t *testing.T) {
	configFile := config.Cfg
	e := NewEmail(configFile.Email)
	e.RunLog("error when: %s", "send abc")
}
