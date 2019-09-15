package app

import (
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"
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

	app.handlers = api.NewHandlers()

	app.handlers.OnStart()

	router := gin.Default()

	v0 := router.Group("/v0")
	{
		v0.GET("/health", app.handlers.HandlerGet)
	}

	v1 := router.Group("/v1")
	{
		v1.GET("/ws/server", app.handlers.HandlerGet)
	}

	fmt.Println("Listen:", g_ConfigData.Port)

	http.ListenAndServe(":"+g_ConfigData.Port, router)

	return nil
}

func (app *App) Shutdown() {
	app.handlers.OnClose()
	fmt.Println("server shutdown")
}
