package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongyuefan/tmpserver/ethscan"
	"github.com/hongyuefan/tmpserver/models"
	"github.com/hongyuefan/tmpserver/types"
	"github.com/hongyuefan/tmpserver/util/log"
)

type Handlers struct {
}

func NewHandlers() *Handlers {

	return &Handlers{}
}

func (h *Handlers) OnClose() {

}

func (h *Handlers) HandlerAddMoney(c *gin.Context) {
	var (
		err    error
		reqAdd types.ReqAddMoney
		member models.Member
		record models.AddMoneyRecord
		block  int64
	)

	if err = c.BindJSON(&reqAdd); err != nil {
		goto errDeal
	}

	member.UserName = reqAdd.UserName
	member.UID = reqAdd.UserId

	if err = models.GetMember(&member, "uid", "username"); err != nil {
		goto errDeal
	}

	record.Address = reqAdd.Address
	record.Hash = reqAdd.Hash
	record.Money = reqAdd.Amount
	record.Status = types.STATUS_WAITING
	record.Type = reqAdd.Type
	record.UID = member.UID
	record.Time = reqAdd.Time

	if block, err = ethscan.GetLastBlock(); err != nil {
		log.GetLog().LogWrite("GetLastBlock Error,Uid:", record.UID, "Address:", record.Address, "Hash:", record.Hash, "Money:", record.Money)
	} else {
		record.CheckedBlock = block - 2
	}

	if _, err = models.AddRecord(&record); err != nil {
		goto errDeal
	}
	HandleSuccessMsg(c, "HandlerAddMoney", "success")
	return
errDeal:
	HandleErrorMsg(c, "HandlerAddMoney", err.Error())
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
