package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"imagecut/internal/config"
	"imagecut/internal/img"
	"imagecut/internal/lru"
	"net/http"
	"sync"
)

type Api struct {
	sync.Mutex
	imgService *img.Img
	cache      *lru.Lru
	cachePath  string
	logOnErr   func(ctx *gin.Context, err error)
}

func NewApi(cacheSize uint, cachePath string, imgConfig config.Img, errorLogger *zap.Logger) *Api {
	logOnError := makeLogOnErr(errorLogger)

	api := &Api{
		logOnErr:   logOnError,
		imgService: img.NewImg(imgConfig.ImageFolder, imgConfig.DownloadTimeout),
		cache:      lru.NewLru(cacheSize, cachePath),
		cachePath:  cachePath,
	}

	err := api.restoreCache()
	errorLogger.Warn("Failed to restore cache", zap.Error(err))
	return api
}

func (a *Api) Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}

func makeLogOnErr(logger *zap.Logger) func(ctx *gin.Context, err error) {
	return func(ctx *gin.Context, err error) {
		var message string

		if ctx != nil {
			message = fmt.Sprintf("\nmethod: %s\n uri: %s\n host: %s\n error:",
				ctx.Request.Method,
				ctx.Request.RequestURI,
				ctx.Request.Host,
			)
		}

		if err != nil {
			logger.Error(message, zap.Error(err))
		}
	}
}
