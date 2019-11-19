package api

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"imagecut/internal/img"
	"net/http"
	"strconv"
)

type Cache interface {
	Set(key string, value interface{}, size uint) ([]interface{}, error)
	Get(key string) (interface{}, error)
}

type Api struct {
	imgService *img.Img
	cache Cache
}


func NewApi(cache Cache, imageFolder string) *Api {
	return &Api{
		imgService: img.NewImg(imageFolder),
		cache: cache,
	}
}

func (a *Api) Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}

func (a *Api) Crop(ctx *gin.Context) {
	url := ctx.Query("origin")

	width, height, err := convertCropParams(ctx.Param("width"), ctx.Param("height"))

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	fileName := hasher(ctx.Request.URL.Path)

	imgData, ok := a.getFromCache(fileName)

	if !ok {
		imgData, err = a.imgService.CropByUrl(url, fileName, width, height)

		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		a.setToCache(fileName, imgData)
	}

	ctx.File(imgData.Path)
}

func (a *Api) getFromCache(key string) (img.ImageData, bool) {
	v, _ := a.cache.Get(key)
	//TODO Add logging for get from cache error
	if v != nil {
		return v.(img.ImageData), true
	}

	return img.ImageData{}, false
}

func (a *Api) setToCache(key string, data img.ImageData) {

	//TODO remove excluded files
	_, _ = a.cache.Set(key, data, 1)
}

func convertCropParams(w, h string) ( width int, height int, err error) {
	width, err = strconv.Atoi(w)

	if err != nil {
		return
	}

	height, err = strconv.Atoi(h)

	return
}

func hasher(s string) string {
	return	fmt.Sprintf("%x", md5.Sum([]byte(s)))
}