package api

import (
	"container/list"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hongyuefan/tmpserver/api/userspace"
	"github.com/hongyuefan/tmpserver/util/log"
)

type Handlers struct {
	upGrader  websocket.Upgrader
	chanClose chan struct{}
	handler   *Handler
	gPool     *GPool
	mapTask   map[string]*list.List
}

func NewHandlers() *Handlers {
	return &Handlers{
		upGrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		gPool:     NewGPool(100000, 100),
		handler:   NewHandler(),
		chanClose: make(chan struct{}),
		mapTask:   make(map[string]*list.List, 0),
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

func (h *Handlers) registHandler() {
	h.handler.RegistHandler(TYPE_PING_PONG, h.pong)
	h.handler.RegistHandler(TYPE_USER_LOGIN, userspace.UserLogin)

	h.handler.RegistHandler(TYPE_USER_SEARCH_SCENE, userspace.GetUserSearchScene)
	h.handler.RegistHandler(TYPE_USER_DEFEN_SCENE, userspace.GetUserDefenScene)
	h.handler.RegistHandler(TYPE_USER_QUEEN_SCENE, userspace.GetUserQueenScene)

	h.handler.RegistHandler(TYPE_USER_SEARCH_UPDATE, userspace.UpdateUserSearchAnt)
	h.handler.RegistHandler(TYPE_USER_DEFEN_UPDATE, userspace.UpdateUserDefenAnt)
	h.handler.RegistHandler(TYPE_USER_QUEEN_UPDATE, userspace.UpdateUserQueenAnt)
}
