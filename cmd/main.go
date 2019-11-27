package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"imagecut/api"
	"imagecut/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error

	conf := config.GetConfig()
	env := config.GetEnv()
	logger := makeErrorLogger(env, conf.Logging.ErrorLog)
	a := api.NewApi(conf.CacheSize, conf.CachePath, conf.Img, logger)

	server := &http.Server{
		Addr:    conf.Http.Addr,
		Handler: makeHandler(a, env, conf),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("", zap.Error(err))
			stop <- os.Interrupt
		}
	}()

	<-stop

	if err = a.FlushCache(); err != nil {
		logger.Error("", zap.Error(err))
	}

	if err = graceful(server, 5*time.Second); err != nil {
		logger.Error("", zap.Error(err))
	}

}

func makeHandler(a *api.Api, env string, conf config.Config) *gin.Engine {
	if env == config.EnvType.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(makeLogMiddleware(env, conf.Logging.AccessLog))

	handler.GET("/status", a.Status)
	handler.GET("/crop/:width/:height/", a.Crop)

	return handler
}

func makeErrorLogger(env string, param config.LogParams) *zap.Logger {
	var err error
	var logger *zap.Logger

	switch env {
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
	if env == config.EnvType.Prod {
		gin.DefaultWriter = &lumberjack.Logger{
			Filename:   params.FileName,
			MaxBackups: params.MaxBackups,
			MaxAge:     params.MaxAge,
		}
	}

	return gin.Logger()
}

func graceful(hs *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return hs.Shutdown(ctx)
}
