package api

import (
	"encoding/json"
	"errors"
)

type CallFunc func([]byte) (interface{}, error)

type Handler struct {
	mHander map[int]CallFunc
}

func NewHandler() *Handler {
	return &Handler{
		mHander: make(map[int]CallFunc, 64),
	}
}

func (h *Handler) RegistHandler(typ int, function CallFunc) {
	h.mHander[typ] = function
}

func (h *Handler) HandleMessage(msg []byte) []byte {
	m := new(ReqMessage)
	if err := json.Unmarshal(msg, m); err != nil {
		return h.responseMessage(false, 0, errors.New("msg format error"))
	}
	callFunc, ok := h.mHander[m.Typ]
	if !ok {
		return h.responseMessage(false, m.Typ, errors.New("msg type not find"))
	}
	rspMsg, err := callFunc(m.Data)
	if err != nil {
		return h.responseMessage(false, m.Typ, err)
	}
	return h.responseMessage(true, m.Typ, rspMsg)
}

func (h *Handler) responseMessage(success bool, typ int, msg interface{}) []byte {
	bytMsg, _ := json.Marshal(msg)
	rsp := &RspMessage{Success: success, Typ: typ, Data: bytMsg}
	rspByt, _ := json.Marshal(rsp)
	return rspByt
}
