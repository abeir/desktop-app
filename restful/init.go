package restful

import (
	"github.com/abeir/desktop-app/core/config"
	"github.com/abeir/desktop-app/core/log"
)

var Gobal *gobalContent

type gobalContent struct {
	Application config.ApplicationConfig
	Api config.ApiConfig
}

//TODO 暂时注释掉
//func init(){
//	Gobal = &gobalContent{}
//
//	initApplicationConfig()
//	initApiConfig()
//	initLog(&Gobal.Application)
//}

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