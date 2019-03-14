package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongyuefan/tmpserver/automan"
	"github.com/hongyuefan/tmpserver/ethscan"
	"github.com/hongyuefan/tmpserver/models"
	"github.com/hongyuefan/tmpserver/types"
	"github.com/hongyuefan/tmpserver/util/log"
)

type Handlers struct {
	auto *automan.AutoMan
}

func NewHandlers() *Handlers {

	return &Handlers{}
}

func (h *Handlers) OnClose() {

}

func (h *Handlers) HandlerGetKuaiDi(c *gin.Context) {
	var (
		err     error
		rspAuto types.RspAuto
		info    string
		url     string
	)

	postid := c.Query("postid")

	auto, err := Get("http://www.kuaidi100.com/autonumber/autoComNum?text=" + postid)
	if err != nil {
		goto errDeal
	}

	if err = json.Unmarshal([]byte(auto), &rspAuto); err != nil {
		goto errDeal
	}
	if len(rspAuto.Auto) <= 0 {
		goto errDeal
	}

	url = "http://www.kuaidi100.com/query?type=" + rspAuto.Auto[0].ComCode + "&postid=" + postid

	info, err = Get(url)
	if err != nil {
		goto errDeal
	}

	c.String(200, info)

	return
errDeal:
	HandleErrorMsg(c, "HandlerGetKuaiDi", err.Error())
	return

}

func (h *Handlers) HandlerGetLastBlockNumber(c *gin.Context) {
	var (
		err         error
		blockNumber int64
	)

	if blockNumber, err = ethscan.GetLastBlock(); err != nil {
		goto errDeal
	}

	c.JSON(200, gin.H{
		"isSuccess": true,
		"message":   blockNumber,
	})

	return
errDeal:
	HandleErrorMsg(c, "HandlerGetLastBlockNumber", err.Error())
	return
}

func (h *Handlers) HandlerGetBlock(c *gin.Context) {
	var (
		blockNumber, timeStamp, lucky int64
		blockHash, stimeStamp         string
		base_format                   string = "2006-01-02 15:04:05"
		err                           error
		rspGetBlock                   types.RspBlockHash
	)
	sBlockNumber := c.Query("blocknumber")

	if blockNumber, err = strconv.ParseInt(sBlockNumber, 10, 64); err != nil {
		goto errDeal
	}
	if blockHash, _, stimeStamp, err = ethscan.GetBlockByBlockNumber(blockNumber); err != nil {
		goto errDeal
	}

	timeStamp, _ = strconv.ParseInt(stimeStamp, 0, 64)

	stimeStamp = time.Unix(timeStamp, 0).Format(base_format)

	lucky, err = strconv.ParseInt("0x"+blockHash[len(blockHash)-8:], 0, 64)

	if err != nil {
		goto errDeal
	}

	rspGetBlock.BlockHash = blockHash
	rspGetBlock.Time = stimeStamp
	rspGetBlock.BlockNumber = blockNumber
	rspGetBlock.LuckyNumber = lucky

	c.JSON(200, gin.H{
		"isSuccess": true,
		"message":   rspGetBlock,
	})

	return
errDeal:
	HandleErrorMsg(c, "HandlerGetBlock", err.Error())
	return
}

func (h *Handlers) HandlerAutoManStop(c *gin.Context) {

	if h.auto != nil {
		h.auto.OnClose()
	}
	HandleSuccessMsg(c, "HandlerAutoManStop", "success")
	return
}

func (h *Handlers) HandlerAutoManStart(c *gin.Context) {
	var (
		err      error
		intervel int64
	)
	sTime := c.Query("intervel")

	if intervel, err = strconv.ParseInt(sTime, 10, 64); err != nil {
		goto errDeal
	}

	if h.auto != nil {
		h.auto.OnClose()
	}

	h.auto = automan.NewAutoMan(intervel)

	go h.auto.OnStart()

	HandleSuccessMsg(c, "HandlerAutoManStart", "success")
	return
errDeal:
	HandleErrorMsg(c, "HandlerAutoManStart", err.Error())
	return
}

func (h *Handlers) HandlerAddMembers(c *gin.Context) {
	var (
		err   error
		count int64
		index int
		uid   int64
	)
	sCount := c.Query("count")

	if count, err = strconv.ParseInt(sCount, 10, 64); err != nil {
		goto errDeal
	}

	for index = 0; index < int(count); index++ {

		if uid, err = models.AddMember(&models.Member{
			Email:      GetRandEmail(),
			PassWord:   "0830297d73b08e3aaf42ca9905b30ed1",
			Img:        "photo/member.jpg",
			GroupId:    1,
			EmailCode:  "1",
			MobileCode: "-1",
			PassCode:   "-1",
			Money:      0,
			Level:      3,
			Time:       time.Now().Unix(),
		}); err != nil {
			goto errDeal
		}
		if err = models.UpdateMember(&models.Member{UID: uid, Money: 8888.0}, "money"); err != nil {
			goto errDeal
		}

	}

	HandleSuccessMsg(c, "HandlerAddMembers", "success")
	return
errDeal:
	HandleErrorMsg(c, "HandlerAddMembers", err.Error())
	return
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

	member.Email = reqAdd.UserName
	member.UID = reqAdd.UserId

	if err = models.GetMember(&member, "uid", "email"); err != nil {
		goto errDeal
	}

	record.Address = strings.ToLower(reqAdd.Address)
	record.Hash = strings.ToLower(reqAdd.Hash)
	record.Money = reqAdd.Amount
	record.Status = types.STATUS_WAITING
	record.Type = reqAdd.Type
	record.UID = member.UID
	record.Time = reqAdd.Time

	if record.Hash == "" {
		record.Hash = GenCode(64)
	}

	if block, err = ethscan.GetLastBlock(); err != nil {
		log.GetLog().LogWrite("GetLastBlock Error,Uid:", record.UID, "Address:", record.Address, "Hash:", record.Hash, "Money:", record.Money)
	} else {
		record.CheckedBlock = block - 1
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

func GetRandEmail() string {

	var email string

	expends := []string{"@126.com", "@163.com", "@yahoo.com", "@gmail.com", "@sina.com", "@139.com"}

	cCount := GetRand(8, 16)

	for i := 0; i < int(cCount); i++ {
		email += GetChar_Low()
		time.Sleep(time.Nanosecond)
	}

	email += expends[int(GetRand(0, float64(len(expends)-1)))]

	return email
}

func Get(url string) (b string, err error) {
	body, err := http.Get(url)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, body.Body)
	b = buf.String()
	return
}

func GetChar_Num() (c string) {
	return string(byte(GetRand(48, 58)))
}
func GetChar_Cap() (c string) {
	return string(byte(GetRand(65, 91)))
}
func GetChar_Low() (c string) {
	return string(byte(GetRand(97, 123)))
}
func GetRand(min float64, max float64) (result float64) {
	source := rand.NewSource(time.Now().UnixNano())
	nRand := rand.New(source)
	return nRand.Float64()*(max-min) + min
}
func GenCode(n int) (code string) {
	var (
		rand int
		str  string
	)
	for i := 0; i < n; i++ {
		rand = int(GetRand(0, 3))
		time.Sleep(time.Nanosecond)
		switch rand {
		case 0:
			str += GetChar_Num()
			continue
		case 1:
			str += GetChar_Low()
			continue
		case 2:
			str += GetChar_Cap()
			continue
		}
	}
	return str
}
