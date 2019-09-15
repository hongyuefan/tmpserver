package api

import (
	"encoding/json"
)

func (h *Handlers) registHandler() {
	h.handler.RegistHandler(TYPE_USER_LOGIN, userLogin)
}

func userLogin(reqMsg []byte) (interface{}, error) {
	user := new(ReqUserLogin)
	if err := json.Unmarshal(reqMsg, user); err != nil {
		return nil, err
	}
	return &RspUserLogin{Token: "123"}, nil
}
