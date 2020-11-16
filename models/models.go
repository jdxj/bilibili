package models

import (
	"encoding/json"
	"io"
)

func NewAPIResponse(body io.Reader) (*APIResponse, error) {
	ar := new(APIResponse)
	decoder := json.NewDecoder(body)
	return ar, decoder.Decode(ar)
}

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	TTL     int             `json:"ttl"`
	Status  bool            `json:"status"`
	Data    json.RawMessage `json:"data"`
}

func (ar *APIResponse) LoginInfo() (*LoginInfo, error) {
	li := new(LoginInfo)
	return li, json.Unmarshal(ar.Data, li)
}

func (ar *APIResponse) SignInfo() (*SignInfo, error) {
	si := new(SignInfo)
	return si, json.Unmarshal(ar.Data, si)
}

func (ar *APIResponse) CoinInfo() (*CoinInfo, error) {
	ci := new(CoinInfo)
	return ci, json.Unmarshal(ar.Data, ci)
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
