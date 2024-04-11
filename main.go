package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gin/models"
	"go-gin/pkg/logging"
	"go-gin/pkg/setting"
	"go-gin/pkg/translations"
	"go-gin/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	translations.Setup("zh")
	//gredis.Setup()
	//util.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description go-gin
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.Info("[info] start http server listening %s", endPoint)

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logging.Debug("Server was shutdown gracefully")
			return
		}
		logging.Error("Server error: %v", err)

	}()

	//优雅退出
	gracefulExitWeb(server)

}

// gracefulExitWeb 优雅退出
func gracefulExitWeb(server *http.Server) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch

	logging.Debug("got a signal %v", sig)
	now := time.Now()
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(cxt)
	if err != nil {
		logging.Debug("%v", err)
	}

	//关闭数据库
	models.CloseDB()

	// 看看实际退出所耗费的时间
	logging.Debug("------exited------  %v", time.Since(now))

	//日志落盘
	if err := logging.GetLogger().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		logging.Error("%v", err)
	}
}
