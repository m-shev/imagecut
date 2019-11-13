package api

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

func Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}

func Crop(ctx *gin.Context) {
	width := ctx.Param("width")
	height := ctx.Param("height")
	url := ctx.Param("url")
	fmt.Println("before replace", url)
	url = strings.Replace(url, "/", "", 1)
	url = strings.Replace(url, "/", "//", 1)
	fmt.Println("after replace", url)
	path := "some"
	err := DownloadFile(url, path)
	fmt.Println("++++++++++", err)

	ctx.String(http.StatusOK, "%s %s %s", width, height, url)
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("download error", err)
		return err
	}
	defer resp.Body.Close()
	img, err := imaging.Decode(resp.Body)
	fmt.Println(err, img)
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}