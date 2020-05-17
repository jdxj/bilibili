package modules

import "encoding/json"

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	TTL     int             `json:"ttl"`
	Status  bool            `json:"status"`
	Data    json.RawMessage `json:"data"`
}

type LoginInfo struct {
	IsLogin bool   `json:"isLogin"`
	Money   int    `json:"money"`
	Uname   string `json:"uname"`
}

type SignInfo struct {
	List  []*SignEntry `json:"list"`
	Count int          `json:"count"`
}

type SignEntry struct {
	Time   string `json:"time"`
	Delta  int    `json:"delta"`
	Reason string `json:"reason"`
}

type CoinInfo struct {
	Money int `json:"money"`
}
