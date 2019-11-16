package api

import (
	"github.com/gin-gonic/gin"
	"imagecut/internal/img"
	"net/http"
	"strconv"
)

type Api struct {
	imgService *img.Img
}

func NewApi() *Api {
	return &Api{
		imgService: img.NewImg(),
	}
}

func (api *Api) Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}

func (api *Api) Crop(ctx *gin.Context) {
	url := ctx.Query("origin")

	width, height, err := convertCropParams(ctx.Param("width"), ctx.Param("height"))

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	image, err := api.imgService.CropFromUrl(url, width, height)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.File(image.Path)
}

func convertCropParams(w, h string) (int, int, error) {
	width, err := strconv.Atoi(w)

	if err != nil {
		return 0, 0, err
	}

	height, err := strconv.Atoi(h)

	if err != nil {
		return 0, 0, err
	}

	return width, height, err
}

