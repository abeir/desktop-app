package restful

import (
	"context"
	"github.com/abeir/desktop-app/core"
	"github.com/abeir/desktop-app/core/log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServerState 服务状态
type ServerState int

const (
	// ServerStarting 服务启动中
	ServerStarting ServerState	 = 1
	// ServerStarted 服务已启动
	ServerStarted ServerState	 = 2
	// ServerFailed 服务启动失败
	ServerFailed ServerState	 = 3
)

// ServerStartedListener 服务启动后调用的监听程序
type ServerStartedListener func(serv *Server)

// NewServer 创建服务
func NewServer() *Server{
	return &Server{
		startedChan: make(chan bool, 1),
		testPortDelayed: time.Second * 2,
		testProtRetryTimes: 3,
	}
}


type Server struct {
	startedChan chan bool
	//延时测试端口的时间
	testPortDelayed time.Duration
	//测试端口的重试次数，若设置为小于1的数则按1次处理
	testProtRetryTimes int
	//服务启动后调用的监听程序
	serverStartedListener ServerStartedListener

	// 端口号
	Port string
	// 服务状态
	State ServerState
}

// ServerStartedListener 设置服务启动后的监听程序
func (s *Server) ServerStartedListener(listener ServerStartedListener) *Server{
	s.serverStartedListener = listener
	return s
}

// listenServerStarted 启动监听，在服务启动时检测http服务监听的端口
func (s *Server) listenServerStarted(){
	if s.serverStartedListener == nil {
		return
	}
	go func(){
		<- s.startedChan
		if s.testPortRetry() {
			s.State = ServerStarted
			log.Info("Server started.")
			s.serverStartedListener(s)
		}
	}()
}

// testPortRetry 检测http服务监听的端口，该方法会延时阻塞执行，若监测端口超时则会重试
func (s *Server) testPortRetry() bool{
	time.Sleep(s.testPortDelayed)
	testProtTimes := core.If(s.testProtRetryTimes < 1, 1, s.testProtRetryTimes).(int)
	for i:=0; i<testProtTimes; i++ {
		if s.testPort() {
			return true
		}
	}
	return false
}

// testPort 检测一次http服务监听的端口
func (s *Server) testPort() bool{
	conn, err := net.DialTimeout("tcp", ":" + s.Port, time.Millisecond * 500)
	defer core.CloseQuietly(conn)
	return err==nil
}

// Start 启动服务
func (s *Server) Start(){
	s.State = ServerStarting
	app := Gobal.Application
	engine := Gobal.engine

	s.Port = app.Server.Port
	serv := &http.Server{Addr:":" + s.Port, Handler: engine}

	s.startServer(serv)
	s.listenServerStarted()
	s.gracefulShutdown(serv)
}

func (s *Server) startServer(serv *http.Server){
	go func(){
		s.startedChan <- true
		if err := serv.ListenAndServe(); err!=nil && err!=http.ErrServerClosed {
			s.State = ServerFailed
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// gracefulShutdown 优雅关闭服务
func (s *Server) gracefulShutdown(serv *http.Server){
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