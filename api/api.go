package api

import (
	"github.com/gin-gonic/gin"
	"imagecut/internal/img"
	"imagecut/internal/lru"
	"log"
	"net/http"
	"sync"
)

type Api struct {
	sync.Mutex
	imgService *img.Img
	cache      *lru.Lru
	cachePath  string
}

func NewApi(cacheSize uint, cachePath string, imageFolder string) *Api {
	api := &Api{
		imgService: img.NewImg(imageFolder),
		cache:      lru.NewLru(cacheSize, cachePath),
		cachePath:  cachePath,
	}

	err := api.restoreCache()
	log.Println(err)
	return api
}

func (a *Api) Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}

func (a *Api) Graceful() error {
	return a.flushCache()
}
