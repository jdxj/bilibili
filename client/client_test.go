package client

import (
	"fmt"
	"testing"
	"time"

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
	cookieCfg := config.GetCookies()
	client := NewClient()

	client.changeCookies(cookieCfg[0].Values)

	if !client.verifyLogin() {
		t.Fatalf("login faild")
	}
	client.sign()
}

func TestClient_Start2(t *testing.T) {
	client := NewClient()
	client.Start()
}

func TestParseTime(t *testing.T) {
	value := "2020-05-17 08:00:02"
	dateTime, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Println(dateTime.Date())
	now := time.Now()

	if now.Year() != dateTime.Year() &&
		now.Month() != dateTime.Month() &&
		now.Day() != dateTime.Day() {
		t.Fatalf("not equal")
	}
}

func TestClientCheckSignLog(t *testing.T) {
	cfg := config.GetCookies()
	client := NewClient()
	client.changeCookies(cfg[0].Values)
	result := client.alreadySign()
	fmt.Printf("%v", result)
}

func TestClient_MulAlreadySign(t *testing.T) {
	cfg := config.GetCookies()
	client := NewClient()
	client.changeCookies(cfg[0].Values)
	result := client.mulAlreadySign()
	if result {
		client.sendCoinNum()
	}
}

func TestClientSendCoinNum(t *testing.T) {
	notifier.AddRecipients("985759262@qq.com")
	cfg := config.GetCookies()
	client := NewClient()
	client.changeCookies(cfg[0].Values)
	client.sendCoinNum()
}
