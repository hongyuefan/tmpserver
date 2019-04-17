package app

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/orm"
	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hongyuefan/tmpserver/api"
	"github.com/hongyuefan/tmpserver/util/config"
	"github.com/hongyuefan/tmpserver/util/log"
)

const MasterName = "myserver"

type ConfigData struct {
	Port    string
	Idls    float64
	LogDir  string
	SqlConn string
	Key     string
	Scr     string
}

type App struct {
	handlers *api.Handlers
}

var g_ConfigData *ConfigData

func OnInitFlag(c *config.Config) (err error) {

	g_ConfigData = new(ConfigData)
	g_ConfigData.Port = c.GetString("port")
	g_ConfigData.Idls = c.GetFloat("idls")
	g_ConfigData.LogDir = c.GetString("logdir")
	g_ConfigData.SqlConn = c.GetString("sqlconn")
	g_ConfigData.Key = c.GetString("key")
	g_ConfigData.Scr = c.GetString("scr")

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

	api.OnInit()

	app.handlers = api.NewHandlers(g_ConfigData.Key, g_ConfigData.Scr)

	router := gin.Default()

	v0 := router.Group("/v0")
	{
		v0.GET("/health", app.handlers.HandlerGet)
	}

	v1 := router.Group("/v1")
	{
		v1.POST("/post", app.handlers.HandlerPost)
		v1.GET("/get", app.handlers.HandlerGet)
		v1.POST("/face/member/add", app.handlers.HandlerAddMember)
		v1.POST("face/result/post", app.handlers.HandlerDetection)
		v1.GET("face/result/get", app.handlers.HandlerDetection)
	}

	fmt.Println("Listen:", g_ConfigData.Port)

	http.ListenAndServe(":"+g_ConfigData.Port, router)

	return nil
}

func (app *App) Shutdown() {
	app.handlers.OnClose()
	fmt.Println("server shutdown")
}
