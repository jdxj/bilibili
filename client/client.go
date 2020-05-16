package client

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/jdxj/bilibili/modules"

	"github.com/jdxj/bilibili/config"
	"github.com/jdxj/bilibili/email"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36"
const Domain = ".bilibili.com"

func NewClient(config *config.Cookie) *Client {
	URL, err := url.Parse(config.Domain)
	if err != nil {
		email.Log("new client-url parse error: %s", err)
		panic(err)
	}

	client := &Client{
		config:  config,
		url:     URL,
		hClient: &http.Client{},
	}
	return client
}

type Client struct {
	config  *config.Cookie
	url     *url.URL
	hClient *http.Client
}

func (c *Client) Start() {
	cfg := c.config

	for _, v := range cfg.Cookies {
		// 1
		c.changeCookies(v)
		// 2
		money := c.verifyLogin()
		if money < 0 {
			continue
		}
		// 3
		c.sign()
		email.Log("sign ok, money: %d", money)
	}
}

func (c *Client) changeCookies(cookie string) {
	hc := c.hClient
	hc.Jar = nil // 清除上一个 cookies

	// 解析新 cookies 并生成新 jar
	cookies := parseCookies(cookie)
	for _, v := range cookies {
		v.Domain = Domain
	}
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(c.url, cookies)

	hc.Jar = jar
}

// 1. 验证登陆
func (c *Client) verifyLogin() int {
	hc := c.hClient

	api := "https://api.bilibili.com/x/web-interface/nav"
	req := newRequestUserAgent(api)
	resp, err := hc.Do(req)
	if err != nil {
		email.Log("verifyLogin-http client do error: %s", err)
		return -1
	}
	defer resp.Body.Close()

	apiResp := &modules.APIResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(apiResp); err != nil {
		email.Log("verifyLogin-api resp decode error: %s", err)
		return -1
	}

	if apiResp.TTL == 0 {
		email.Log("verifyLogin-api resp ttl error: %d", apiResp.TTL)
		return -1
	}

	loginInfo := &modules.LoginInfo{}
	if err := json.Unmarshal(apiResp.Data, loginInfo); err != nil {
		email.Log("verifyLogin-login info unmarshal error: %d", err)
		return -1
	}
	return loginInfo.Money
}

// 2. 签到
func (c *Client) sign() {
	hc := c.hClient
	cfg := c.config

	req := newRequestUserAgent(cfg.Domain)
	resp, err := hc.Do(req)
	if err != nil {
		email.Log("sign-http do error: %s", err)
		return
	}
	defer resp.Body.Close()
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
