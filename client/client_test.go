package client

import (
	"fmt"
	"testing"

	"github.com/jdxj/bilibili/config"
)

func TestParseCookies(t *testing.T) {
	cookies := "a=b;c=d"
	result := parseCookies(cookies)
	for _, cookie := range result {
		fmt.Printf("%#v\n", *cookie)
	}
}

func TestClient_Start(t *testing.T) {
	cookieCfg := config.Cfg.Cookie
	client := NewClient(cookieCfg)

	client.changeCookies(cookieCfg.Cookies[0])

	if client.verifyLogin() < 0 {
		t.Fatalf("login faild")
	}
	client.sign()
}
