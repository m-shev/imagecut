package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (a *Api) Crop(ctx *gin.Context) {
	url := ctx.Query("origin")

	width, height, err := convertCropParams(ctx.Param("width"), ctx.Param("height"))

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	fileName := hasher(ctx.Request.RequestURI)

	imgData, ok := a.getFromCache(fileName, ctx)

	if !ok {
		imgData, err = a.imgService.CropByUrl(url, fileName, width, height)

		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		a.setToCache(fileName, imgData, ctx)
	}

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
