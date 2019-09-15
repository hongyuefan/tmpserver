package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hongyuefan/tmpserver/util/log"
)

type Handlers struct {
	upGrader  websocket.Upgrader
	chanClose chan struct{}
	handler   *Handler
	gPool     *GPool
}

func NewHandlers() *Handlers {
	return &Handlers{
		upGrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		gPool:     NewGPool(1, 1),
		handler:   NewHandler(),
		chanClose: make(chan struct{}),
	}
}

func (h *Handlers) OnStart() {
	h.registHandler()
	h.gPool.Start()
}

func (h *Handlers) OnClose() {
	close(h.chanClose)
	h.gPool.Stop()
}

func (h *Handlers) HandlerGet(c *gin.Context) {
	ws, err := h.upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		select {
		case <-h.chanClose:
			return
		default:
			if err := h.handlerMessage(ws); err != nil {
				log.GetLog().LogError("ws connect", ws.RemoteAddr().String(), "error", err)
				return
			}
		}
	}
}

func (h *Handlers) handlerMessage(ws *websocket.Conn) error {
	typ, msg, err := ws.ReadMessage()
	if err != nil {
		return err
	}
	rsp := h.handler.HandleMessage(msg)
	h.gPool.SendMsg(typ, rsp, ws)
	return nil
}
