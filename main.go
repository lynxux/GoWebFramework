package main

import (
	"context"
	"github.com/lynxux/goWebFramework/framework/gin"
	"github.com/lynxux/goWebFramework/provider/demo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core := gin.New()

	core.Bind(&demo.DemoServiceProvider{})

	core.Use(gin.Recovery())
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8088",
	}
	go func() {
		server.ListenAndServe()
	}()

	// 因为使用 Ctrl 或者 kill 命令，它们发送的信号是进入 main 函数的，即只有 main 函数所在的 Goroutine 会接收到，所以必须在 main 函数所在的 Goroutine 监听信号
	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 超时的context
	defer cancel()

	//Shutdown() 一旦执行之后，它会阻塞当前 Goroutine，并且在所有连接请求都结束之后，才继续往后执行
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown", err)
	}

}
