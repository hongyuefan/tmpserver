package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongyuefan/tmpserver/types"
	"github.com/hongyuefan/tmpserver/util/log"
)

type Handlers struct {
	Key string
	Scr string
}

func NewHandlers(key, scr string) *Handlers {

	return &Handlers{
		Key: key,
		Scr: scr,
	}
}

func (h *Handlers) OnClose() {

}

func (h *Handlers) HandlerDetection(c *gin.Context) {
	var (
		err    error
		result string
	)
	url := c.Query("url")

	if result, err = GetResult(h.Key, h.Scr, url); err != nil {
		goto errDeal
	}

	HandleSuccessMsg(c, "HandlerDetection", result)
	return
errDeal:
	HandleErrorMsg(c, "HandlerAddMember", err.Error())
	return
}

func (h *Handlers) HandlerAddMember(c *gin.Context) {
	var (
		err    error
		reqAdd types.ReqMember
	)

	if err = c.BindJSON(&reqAdd); err != nil {
		goto errDeal
	}

	HandleSuccessMsg(c, "HandlerAddMember", "success")
	return
errDeal:
	HandleErrorMsg(c, "HandlerAddMember", err.Error())
	return
}

func (h *Handlers) HandlerPost(c *gin.Context) {
	var (
		err error
	)
	if err != nil {
		goto errDeal
	}
	return
errDeal:
	HandleErrorMsg(c, "HandlerPost", err.Error())
	return
}

func (h *Handlers) HandlerGet(c *gin.Context) {
	var (
		err error
	)
	if err != nil {
		goto errDeal
	}
	return
errDeal:
	HandleErrorMsg(c, "HandlerGet", err.Error())
	return
}

func HandleSuccessMsg(c *gin.Context, requestType, msg string) {
	responseWrite(c, true, msg)
	logMsg := fmt.Sprintf("type[%s] From [%s] Params [%s]", requestType, c.Request.RemoteAddr, msg)
	log.GetLog().LogInfo(logMsg)
}

func HandleDebugMsg(c *gin.Context, requestType string, info string) {
	logMsg := fmt.Sprintf("type[%s] From [%s] Params [%s]", requestType, c.Request.RemoteAddr, info)
	log.GetLog().LogDebug(logMsg)
}
func HandleErrorMsg(c *gin.Context, requestType string, result string) {
	msg := fmt.Sprintf("type[%s] From [%s] Error [%s] ", requestType, c.Request.RemoteAddr, result)
	responseWrite(c, false, msg)
	log.GetLog().LogError(msg)
}
func responseWrite(ctx *gin.Context, isSuccess bool, result string) {
	ctx.JSON(200, gin.H{
		"isSuccess": isSuccess,
		"message":   result,
	})
}
