package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/jdxj/bilibili/config"

	"github.com/jdxj/bilibili/modules"

	"github.com/jdxj/bilibili/email"
)

const (
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36"

	Domain   = ".bilibili.com"
	LoginAPI = "https://api.bilibili.com/x/web-interface/nav"
	SignAPI  = "https://api.bilibili.com/x/member/web/coin/log?jsonp=jsonp"
	CoinAPI  = "https://account.bilibili.com/site/getCoin"
	WebSite  = "https://www.bilibili.com"

	TimeLayout = "2006-01-02 15:04:05"
)

var (
	notifier *email.Email
)

func init() {
	notifier = email.NewEmail()
}

func NewClient() *Client {
	target, _ := url.Parse(WebSite)

	client := &Client{
		url:     target,
		hClient: &http.Client{},
	}
	return client
}

type Client struct {
	url     *url.URL
	hClient *http.Client

	jars []http.CookieJar
}

func (c *Client) Start() {
	cookiesConfig := config.GetCookies()
	for _, cookie := range cookiesConfig {
		notifier.ResetRecipients()
		notifier.AddRecipients(cookie.Recipient)

		c.changeCookies(cookie.Values)
		if !c.verifyLogin() {
			continue
		}
		c.sign()

		if !c.mulAlreadySign() {
			continue
		}
		c.sendCoinNum()
	}
}

// 0. 装载 cookie
func (c *Client) changeCookies(cookie string) {
	hc := c.hClient

	// 解析新 cookies 并生成新 jar
	cookies := parseCookies(cookie)
	for _, v := range cookies {
		v.Domain = Domain
	}
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(c.url, cookies)

	hc.Jar = jar
	c.jars = append(c.jars, jar)
}

// 1. 验证登陆
func (c *Client) verifyLogin() bool {
	hc := c.hClient

	req := newRequestUserAgent(LoginAPI)
	resp, err := hc.Do(req)
	if err != nil {
		email.Log("verifyLogin-http client do error: %s", err)
		return false
	}
	defer resp.Body.Close()

	apiResp := &modules.APIResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(apiResp); err != nil {
		email.Log("verifyLogin-api resp decode error: %s", err)
		return false
	}

	if apiResp.TTL == 0 {
		email.Log("verifyLogin-api resp ttl error: %d", apiResp.TTL)
		return false
	}

	loginInfo := &modules.LoginInfo{}
	if err := json.Unmarshal(apiResp.Data, loginInfo); err != nil {
		email.Log("verifyLogin-login info unmarshal error: %d", err)
		return false
	}
	return true
}

// 2. 签到
func (c *Client) sign() {
	hc := c.hClient

	req := newRequestUserAgent(WebSite)
	resp, err := hc.Do(req)
	if err != nil {
		email.Log("sign-http do error: %s", err)
		return
	}
	resp.Body.Close()
}

func (c *Client) mulAlreadySign() bool {
	num := 10                                 // 重试10次
	ticker := time.NewTicker(5 * time.Second) // 间隔 5s
	defer ticker.Stop()

	for {
		<-ticker.C
		if c.alreadySign() {
			return true
		}

		num--
		if num <= 0 {
			break
		}
	}

	format := "已经执行了签到程序并进行了硬币数量检测, 但仍未检测到硬币更新, 可能是B站还未统计. 请手动查看硬币获取记录: %s"
	addr := "https://account.bilibili.com/account/coin"
	notifier.SignLog(format, addr)
	email.Log("sign fail, addr: %v", notifier.To())
	return false
}

// 3. 检查是否已获得
func (c *Client) alreadySign() bool {
	hc := c.hClient

	req := newRequestUserAgent(SignAPI)
	resp, err := hc.Do(req)
	if err != nil {
		email.Log("alreadySign-hc.Do error: %s", err)
		return false
	}
	defer resp.Body.Close()

	apiResp := &modules.APIResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(apiResp); err != nil {
		email.Log("alreadySign-decoder.Decode1 error: %s", err)
		return false
	}

	if apiResp.TTL == 0 {
		email.Log("alreadySign-api resp ttl error: %d", apiResp.TTL)
		return false
	}

	signInfo := &modules.SignInfo{}
	if err := json.Unmarshal(apiResp.Data, signInfo); err != nil {
		email.Log("alreadySign-unmarshal signInfo error: %s", err)
		return false
	}

	if len(signInfo.List) <= 0 {
		email.Log("alreadySign-can not get sign log: %d", signInfo.Count)
		return false
	}

	signEntry := signInfo.List[0]
	curDate, _ := time.Parse(TimeLayout, signEntry.Time)
	now := time.Now()

	if curDate.Year() != now.Year() &&
		curDate.Month() != now.Month() &&
		curDate.Day() != now.Day() {

		return false
	}
	return true
}

// 4. 保存 cookie, 可能会有用
func (c *Client) SaveCookies() {
	for i, v := range c.jars {
		cookies := v.Cookies(c.url)
		fmt.Printf("%02d--------", i)
		for _, vv := range cookies {
			fmt.Printf("%s\n", vv)
		}
	}
}

func (c *Client) sendCoinNum() {
	hc := c.hClient

	req := newRequestUserAgent(CoinAPI)
	resp, err := hc.Do(req)
	if err != nil {
		email.Log("sendCoinNum-hc.Do error: %s", err)
		return
	}
	defer resp.Body.Close()

	apiResp := &modules.APIResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(apiResp); err != nil {
		email.Log("sendCoinNum-decoder.Decode1 error: %s", err)
		return
	}

	if !apiResp.Status {
		email.Log("sendCoinNum-api resp status error: %v", apiResp.Status)
		return
	}

	coinInfo := &modules.CoinInfo{}
	if err := json.Unmarshal(apiResp.Data, coinInfo); err != nil {
		email.Log("sendCoinNum-unmarshal coinInfo error: %s", err)
		return
	}

	notifier.SignLog("sign success, money: %d", coinInfo.Money)
}

func parseCookies(cookies string) []*http.Cookie {
	header := http.Header{}
	header.Set("Cookie", cookies)

	req := &http.Request{Header: header}
	return req.Cookies()
}

func newRequestUserAgent(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", UserAgent)
	return req
}
