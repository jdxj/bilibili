package config

import (
	"encoding/json"
	"os"
)

/*
config file
{
  "email": {
    "user": "985759262@qq.com",
    "password": ""
  },
  "cookie": {
    "domain": "https://www.bilibili.com",
    "cookies": []
  }
}

crontab
0 8 * * * cd /root/bilibili && ./bilibili.out
*/

var (
	cfg *Config
)

func init() {
	cfg = readConfig("./config.json")
}

type Config struct {
	Email   *Email    `json:"email"`
	Cookies []*Cookie `json:"cookies"`
}

type Email struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Cookie struct {
	Recipient string `json:"recipient"`
	Values    string `json:"values"`
}

func readConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		panic(err)
	}
	return config
}

func GetEmail() Email {
	return *cfg.Email
}

func GetCookies() []Cookie {
	var result []Cookie
	for _, v := range cfg.Cookies {
		result = append(result, *v)
	}
	return result
}
