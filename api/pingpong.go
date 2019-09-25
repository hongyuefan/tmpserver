package api

import (
	"container/list"
	"encoding/json"
)

//根据uuid获取任务队列，通过pong下发任务
func (h *Handlers) pong(reqMsg []byte) (interface{}, error) {
	ping := new(Ping)
	if err := json.Unmarshal(reqMsg, ping); err != nil {
		return nil, err
	}
	l, ok := h.mapTask[ping.UUID]
	if !ok {
		ls := list.New()
		h.mapTask[ping.UUID] = ls
		return "", nil
	}
	e := l.Front()
	if e == nil {
		return "", nil
	}
	l.Remove(e)
	return e.Value, nil
}

//接受ping数据，根据请求类型处理
func (h *Handlers) ping(typ int, raw json.RawMessage) {
	switch typ {
	default:
	}
}

func (h *Handlers) addTask(uuid string, task *Pong) {
	l, ok := h.mapTask[uuid]
	if ok {
		l.PushBack(task)
	}
}
