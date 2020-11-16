package models

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/jdxj/bilibili/utils"
)

const (
	TimeLayout = "2006-01-02 15:04:05"

	Domain = ".bilibili.com"
	Home   = "https://www.bilibili.com"

	APILogin = "https://api.bilibili.com/x/web-interface/nav"
	APISign  = "https://api.bilibili.com/x/member/web/coin/log?jsonp=jsonp"
	APICoin  = "https://account.bilibili.com/site/getCoin"
)

func NewBiliBili(cookie string) (*BiliBili, error) {
	u, _ := url.Parse(Home)
	cookies := utils.StringToCookies(Domain, cookie)

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)

	client := &http.Client{}
	client.Jar = jar

	b := &BiliBili{
		c: client,
	}

	err := b.login()
	if err != nil {
		return nil, err
	}
	return b, nil
}

type BiliBili struct {
	c *http.Client

	li *LoginInfo
	si *SignInfo
	ci *CoinInfo
}

func (b *BiliBili) Run() error {
	err := b.login()
	if err != nil {
		return err
	}

	err = b.sign()
	if err != nil {
		return err
	}

	return b.tryCheckSign()
}

func (b *BiliBili) login() error {
	req := utils.NewRequestGet(APILogin)
	resp, err := b.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ar, err := NewAPIResponse(resp.Body)
	if err != nil {
		return err
	}
	if ar.Code != 0 {
		return fmt.Errorf("%#v", *ar)
	}

	b.li, err = ar.LoginInfo()
	return err
}

func (b *BiliBili) sign() error {
	req := utils.NewRequestGet(Home)
	resp, err := b.c.Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func (b *BiliBili) tryCheckSign() error {
	err := b.checkSign()
	if err == nil {
		return nil
	}

	dur := time.Second
	timer := time.NewTimer(dur)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := b.checkSign(); err == nil {
				return nil
			}

			dur = dur * 2
			timer.Reset(dur)
		}

		if dur >= 8 {
			return fmt.Errorf("stop retry check sign")
		}
	}
}

func (b *BiliBili) checkSign() error {
	req := utils.NewRequestGet(APISign)
	resp, err := b.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ar, err := NewAPIResponse(resp.Body)
	if err != nil {
		return err
	}
	if ar.Code != 0 {
		return fmt.Errorf("%#v", *ar)
	}

	b.si, err = ar.SignInfo()
	if err != nil {
		return err
	}

	if b.si.Count <= 0 {
		return fmt.Errorf("there has never been sign")
	}

	se := b.si.List[0]
	signDate, _ := time.Parse(TimeLayout, se.Time)
	now := time.Now()
	if signDate.Year() != now.Year() ||
		signDate.Month() != now.Month() ||
		signDate.Day() != now.Day() {
		return fmt.Errorf("sign failed")
	}
	return nil
}

func (b *BiliBili) coins() (int, error) {
	req := utils.NewRequestGet(APICoin)
	resp, err := b.c.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	ar, err := NewAPIResponse(resp.Body)
	if err != nil {
		return -1, err
	}
	if ar.Code != 0 {
		return -1, fmt.Errorf("get coin info failed")
	}

	b.ci, err = ar.CoinInfo()
	if err != nil {
		return -1, err
	}
	return b.ci.Money, nil
}

func (b *BiliBili) Subject() string {
	return "硬币余额"
}

func (b *BiliBili) Content() string {
	coin, err := b.coins()
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%d", coin)
}
