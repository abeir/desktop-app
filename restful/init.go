package restful

import (
	"github.com/abeir/desktop-app/core/config"
	"github.com/abeir/desktop-app/core/log"
	"github.com/abeir/desktop-app/restful/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Gobal *gobalContent

type gobalContent struct {
	Application config.ApplicationConfig
	Api config.ApiConfig

	engine *gin.Engine
}


func init(){
	Gobal = &gobalContent{}

	initApplicationConfig()
	initApiConfig()
	initLog(&Gobal.Application)
	initController(&Gobal.Application)
}

func initApplicationConfig(){
	applicationConfig := config.NewApplicationConfig()
	if err := applicationConfig.Load(); err!=nil {
		panic(err)
	}
	Gobal.Application = *applicationConfig
}

func initApiConfig(){
	apiConfig := config.NewApiConfig()
	if err := apiConfig.Load(); err!=nil {
		panic(err)
	}
	Gobal.Api = *apiConfig
}

func initLog(app *config.ApplicationConfig){
	log.InitLog(app)
}

func initController(app *config.ApplicationConfig){
	controller.SetMode(app)

	Gobal.engine = gin.New()

	engine := Gobal.engine

	engine.LoadHTMLGlob("ui/template/**/*")
	engine.StaticFS("assets", http.Dir("ui/assets"))
	engine.Use(gin.Recovery())
	engine.Use(controller.Logger())

	controller.Validator()
	controller.Router(engine)
}
