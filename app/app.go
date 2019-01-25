package app

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/orm"
	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hongyuefan/tmpserver/api"
	"github.com/hongyuefan/tmpserver/backserver"
	"github.com/hongyuefan/tmpserver/util/config"
	"github.com/hongyuefan/tmpserver/util/log"
)

const MasterName = "myserver"

type ConfigData struct {
	Port     string
	Idls     float64
	LogDir   string
	SqlConn  string
	Intervel int64
	Founder  string
	Blocks   int64
}

type App struct {
	handlers *api.Handlers
	backer   *backserver.Server
}

var g_ConfigData *ConfigData

func OnInitFlag(c *config.Config) (err error) {

	g_ConfigData = new(ConfigData)
	g_ConfigData.Port = c.GetString("port")
	g_ConfigData.Idls = c.GetFloat("idls")
	g_ConfigData.LogDir = c.GetString("logdir")
	g_ConfigData.SqlConn = c.GetString("sqlconn")
	g_ConfigData.Founder = c.GetString("founder")
	g_ConfigData.Intervel = int64(c.GetFloat("interval"))
	g_ConfigData.Blocks = int64(c.GetFloat("blocks"))

	if "" == g_ConfigData.Port || 0 == g_ConfigData.Idls || "" == g_ConfigData.LogDir {
		return fmt.Errorf("config not right")
	}
	return

}

func (app *App) OnStart(c *config.Config) error {

	if err := OnInitFlag(c); err != nil {
		return err
	}

	if _, err := log.NewLog(g_ConfigData.LogDir, MasterName, 0); err != nil {
		return err
	}

	if err := orm.RegisterDataBase("default", "mysql", g_ConfigData.SqlConn); err != nil {
		return err
	}

	app.backer = backserver.NewServer(g_ConfigData.Founder, g_ConfigData.Intervel, g_ConfigData.Blocks)

	go app.backer.Handler()

	app.handlers = api.NewHandlers()

	router := gin.Default()

	v0 := router.Group("/v0")
	{
		v0.GET("/health", app.handlers.HandlerGet)
	}

	v1 := router.Group("/v1")
	{
		v1.POST("/post", app.handlers.HandlerPost)
		v1.GET("/get", app.handlers.HandlerGet)
		v1.POST("/recharge/add/record", app.handlers.HandlerAddMoney)
		v1.GET("/automan/start", app.handlers.HandlerAutoManStart)
		v1.GET("/automan/stop", app.handlers.HandlerAutoManStop)
	}

	fmt.Println("Listen:", g_ConfigData.Port)

	http.ListenAndServe(":"+g_ConfigData.Port, router)

	return nil
}

func (app *App) Shutdown() {
	app.handlers.OnClose()
	app.backer.OnClose()
	fmt.Println("server shutdown")
}
