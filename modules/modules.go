package modules

import "encoding/json"

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	TTL     int             `json:"ttl"`
	Data    json.RawMessage `json:"data"`
}

type LoginInfo struct {
	IsLogin bool   `json:"isLogin"`
	Money   int    `json:"money"`
	Uname   string `json:"uname"`
}
