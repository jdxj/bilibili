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

var Cfg *Config

func init() {
	Cfg = readConfig("./config.json")
}

type Config struct {
	Email  *Email  `json:"email"`
	Cookie *Cookie `json:"cookie"`
}

type Email struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Cookie struct {
	Domain  string   `json:"domain"`
	Cookies []string `json:"cookies"`
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
