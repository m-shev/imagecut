package api

import (
	"github.com/gin-gonic/gin"
	"imagecut/internal/img"
	"net/http"
	"strconv"
)

func (a *Api) Crop(ctx *gin.Context) {
	fileName := hasher(ctx.Request.RequestURI)

	imgData, ok := a.getFromCache(fileName, ctx)

	if ok {
		ctx.Header("X-IMAGECUT-FROM-CACHE", "true")
		ctx.Header("cache-control", "public, max-age=3600")
		ctx.File(imgData.Path)
	} else {
		a.downloadAndCrop(ctx, fileName)
	}
}

func (a *Api) downloadAndCrop(ctx *gin.Context, fileName string) {
	url := ctx.Query("origin")

	width, height, err := convertCropParams(ctx.Param("width"), ctx.Param("height"))

	if err != nil {
		a.logOnErr(ctx, err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	imgData, err := a.imgService.CropByUrl(img.ImageSource{
		Url:      url,
		Headers:  &ctx.Request.Header,
		FileName: fileName,
	}, width, height)

	if err != nil {
		var statusCode int

		if imgData.StatusCode != 0 {
			statusCode = imgData.StatusCode
		} else {
			statusCode = http.StatusInternalServerError
		}

		ctx.String(statusCode, err.Error())
		a.logOnErr(ctx, err)
		return
	}

	a.setToCache(fileName, imgData, ctx)
	ctx.Header("X-IMAGECUT-FROM-CACHE", "false")

	for k, v := range imgData.Header {
		for _, h := range v {
			ctx.Header(k, h)
		}
	}

	ctx.Status(imgData.StatusCode)
	ctx.File(imgData.Path)
}

func convertCropParams(w, h string) (width int, height int, err error) {
	width, err = strconv.Atoi(w)

	if err != nil {
		return
	}

	height, err = strconv.Atoi(h)

	return
}