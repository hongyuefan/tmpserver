package api

import (
	"encoding/json"
)

const (
	TYPE_USER_LOGIN = 10001
)

type ReqMessage struct {
	Typ  int             `json:"type"`
	Data json.RawMessage `json:"data"`
}

type RspMessage struct {
	Success bool            `json:"success"`
	Typ     int             `json:"type"`
	Data    json.RawMessage `json:"data"`
}

type ReqUserLogin struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type RspUserLogin struct {
	Token string `json:"token"`
}
