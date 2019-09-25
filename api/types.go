package api

import (
	"encoding/json"
)

//user request type
const (
	TYPE_PING_PONG          = 10000
	TYPE_USER_LOGIN         = 10001
	TYPE_USER_SEARCH_SCENE  = 10002
	TYPE_USER_DEFEN_SCENE   = 10003
	TYPE_USER_QUEEN_SCENE   = 10004
	TYPE_USER_SEARCH_UPDATE = 10005
	TYPE_USER_DEFEN_UPDATE  = 10006
	TYPE_USER_QUEEN_UPDATE  = 10007
)

//server active send msg type
const (
	TYPE_ADD = 20001
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

type Ping struct {
	UUID string          `json:"uuid"`
	Typ  int             `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Pong struct {
	Typ  int             `json:"type"`
	Data json.RawMessage `json:"data"`
}
