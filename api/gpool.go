package api

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/hongyuefan/tmpserver/util/log"
)

type Msg struct {
	typ int
	msg []byte
	ws  *websocket.Conn
}

type GPool struct {
	routineCount int
	chanMsg      chan *Msg
	chanClose    chan struct{}
}

func NewGPool(msgCount, routineCount int) *GPool {
	return &GPool{
		chanMsg:   make(chan *Msg, msgCount),
		chanClose: make(chan struct{}),
	}
}

func (p *GPool) Start() {
	for i := 0; i < p.routineCount; i++ {
		go p.handler()
	}
}

func (p *GPool) Stop() {
	close(p.chanClose)
}

func (p *GPool) handler() {
	for {
		select {
		case msg := <-p.chanMsg:
			if err := msg.ws.WriteMessage(msg.typ, msg.msg); err != nil {
				log.GetLog().LogError("writeMessage to", msg.ws.RemoteAddr().String(), "error", err.Error())
			}
		case <-p.chanClose:
			return
		}
	}
}

func (p *GPool) SendMsg(typ int, msg []byte, ws *websocket.Conn) {
	timeout := time.NewTicker(time.Second * 30)
	defer timeout.Stop()
	select {
	case p.chanMsg <- &Msg{typ, msg, ws}:
	case <-timeout.C:
		log.GetLog().LogError("writeMessage chan full")
	}
}
