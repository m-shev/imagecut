package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
	env := config.GetEnv()

	logger := makeErrorLogger(env, conf.Logging.ErrorLog)

	a := api.NewApi(conf.CacheSize, conf.CachePath, conf.Img, logger)

	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(makeLogMiddleware(env, conf.Logging.AccessLog))

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

func makeErrorLogger(env string, param config.LogParams) *zap.Logger {
	var err error
	var logger *zap.Logger

	switch env {
	case config.EnvType.QA:
		fallthrough
	case config.EnvType.Prod:
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   param.FileName,
			MaxBackups: param.MaxBackups,
			MaxAge:     param.MaxAge,
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.ErrorLevel,
		)
		logger = zap.New(core)
	default:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatal(err)
	}

	return logger
}

func makeLogMiddleware(env string, params config.LogParams) gin.HandlerFunc {
	if env == config.EnvType.QA || env == config.EnvType.Prod {
		gin.DefaultWriter = &lumberjack.Logger{
			Filename:   params.FileName,
			MaxBackups: params.MaxBackups,
			MaxAge:     params.MaxAge,
		}
	}

	return gin.Logger()
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
