package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"imagecut/api"
	"imagecut/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.AddConfigPath("../config")
	conf := config.GetConfig()
	logger := makeLogger(config.GetEnv())

	a := api.NewApi(conf.CacheSize, conf.CachePath, conf.ImageFolder, logger)
	//gin.DefaultWriter = &lumberjack.Logger{
	//	Filename:   "foo.log",
	//	MaxSize:    500, // megabytes
	//	MaxBackups: 3,
	//	MaxAge:     28, // days
	//}

	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(gin.Logger())

	handler.GET("/status", a.Status)
	handler.GET("/crop/:width/:height/", a.Crop)

	server := &http.Server{
		Addr:    conf.Http.Addr,
		Handler: handler,
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			logger.Error("", zap.Error(err))
		}
	}()

	err := graceful(server, 5*time.Second, makeGracefulCb(a, logger))

	if err != nil {
		logger.Error("", zap.Error(err))
	} else {
		logger.Info("Server stopped successfully")
	}
}

func makeLogger(env string) *zap.Logger {
	var err error
	var logger *zap.Logger

	switch env {
	case config.EnvType.Prod:
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatal(err)
	}

	return logger
}

func makeGracefulCb(a *api.Api, logger *zap.Logger) func() {
	return func() {
		if err := a.Graceful(); err != nil {
			logger.Error("", zap.Error(err))
		}
	}
}

func graceful(hs *http.Server, timeout time.Duration, callback func()) error {
	var err error

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	callback()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := hs.Shutdown(ctx); err != nil {
		return err
	}

	return err
}
