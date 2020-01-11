package restful

import (
	"context"
	"github.com/abeir/desktop-app/core/config"
	"github.com/abeir/desktop-app/core/log"
	"github.com/abeir/desktop-app/restful/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	engine.LoadHTMLGlob("ui/**/*")
	engine.StaticFS("assets", http.Dir("ui/assets"))
	engine.Use(gin.Recovery())
	engine.Use(controller.Logger())

	controller.Validator()
	controller.Router(engine)
}

func StartServer(){
	app := Gobal.Application
	serv := &http.Server{Addr:":" + app.Server.Port, Handler: Gobal.engine}

	go func(){
		if err := serv.ListenAndServe(); err!=nil && err!=http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func(){
		cancel()
	}()
	if err := serv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server exiting")
}