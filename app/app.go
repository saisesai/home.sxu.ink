package app

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/saisesai/home.sxu.ink/config"
	"github.com/saisesai/home.sxu.ink/controller"
	"github.com/saisesai/home.sxu.ink/log"
	"github.com/saisesai/home.sxu.ink/mqtt"
	ginLogrus "github.com/toorop/gin-logrus"
	"mime"
	"net/http"
)

func init() {
	// add mine type
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.WithField("error", err).Panicln("failed to add mine ext type!")
	}
}

type App struct {
	Gin  *gin.Engine
	Mqtt paho.Client
}

func (a *App) Run() {
	go func() {
		if token := a.Mqtt.Connect(); token.Wait() && token.Error() != nil {
			log.WithField("error", token.Error()).Fatalln("failed to run mqtt!")
		}
	}()
	err := a.Gin.Run(config.ListenAddress)
	if err != nil {
		log.WithField("error", err).Fatalln("failed to run gin!")
	}
}

func New() *App {
	return &App{
		Gin:  GinNew(),
		Mqtt: mqtt.Client,
	}
}

func GinNew() *gin.Engine {
	var err error
	engine := gin.New()
	// default config but use logrus
	engine.Use(ginLogrus.Logger(log.Logger), gin.Recovery())
	// set trusted proxy
	err = engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.WithField("error", err).Fatalln("failed to set trusted proxy!")
	}
	// static file
	engine.Use(static.Serve("/", static.LocalFile("public", false)))
	// ping
	engine.GET("/api/ping/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "PONG")
	})
	// user base api
	engine.POST("/api/user/login", controller.HandleUserLogin)
	// user devices api
	engine.GET("/api/user/device/get", controller.HandleUserGetDevices)
	engine.POST("/api/user/device/add", controller.HandleUserAddDevice)
	engine.POST("/api/user/device/del", controller.HandleUserDeleteDevice)
	engine.POST("/api/user/device/mod", controller.HandleUserModifyDevice)
	// device api
	engine.POST("/api/device/info", controller.HandleDeviceInfo)
	engine.POST("/api/device/cmd", controller.HandleDeviceCmd)
	return engine
}
