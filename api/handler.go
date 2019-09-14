package api

import (
	"encoding/json"
	"errors"
)

type Hander struct {
}

func (h *Hander) HandleMessage(msg []byte) []byte {
	m := new(ReqMessage)
	if err := json.Unmarshal(msg, m); err != nil {
		return h.responseMessage(false, 0, errors.New("msg format error"))
	}
	switch m.Typ {
	case TYPE_USER_LOGIN:
		rspMsg, err := h.userLogin(m.Data)
		if err != nil {
			return h.responseMessage(false, m.Typ, err)
		}
		return h.responseMessage(true, m.Typ, rspMsg)

	}
	return h.responseMessage(false, m.Typ, errors.New("msg type not find"))
}

func (h *Hander) responseMessage(success bool, typ int, msg interface{}) []byte {
	bytMsg, _ := json.Marshal(msg)
	rsp := &RspMessage{Success: success, Typ: typ, Data: bytMsg}
	rspByt, _ := json.Marshal(rsp)
	return rspByt

}

func (h *Hander) userLogin(reqMsg []byte) (*RspUserLogin, error) {
	user := new(ReqUserLogin)
	if err := json.Unmarshal(reqMsg, user); err != nil {
		return nil, err
	}
	return &RspUserLogin{Token: "123"}, nil
}
