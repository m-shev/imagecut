package img

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"net/http"
	"strings"
)

type Cache interface {
	Set(key string, value interface{}, size uint) ([]interface{}, error)
	Get(key string) (interface{}, error)
}

type Img struct {
	imageFolder string
	cache       Cache
}

type ImageData struct {
	ImgType string
	Id      string
	Path    string
	Header  http.Header
	image   *image.Image
}

func NewImg(cache Cache) *Img {
	return &Img{cache: cache}
}

func (i *Img) CropByUrl(url string, width, height int) (ImageData, error) {
	data, err := i.downloadImage(url)

	if err != nil {
		return data, err
	}

	cropped := imaging.CropAnchor(*data.image, width, height, imaging.Center)
	data.Path = fmt.Sprintf("some.%s", data.ImgType)
	err = imaging.Save(cropped, data.Path)

	if err != nil {
		return data, err
	}

	return data, nil
}

func (i *Img) getFromCache(key string) (ImageData, bool, error) {
	v, err := i.cache.Get(key)

	if err != nil {
		return ImageData{}, false, err
	}

	if v != nil {
		return v.(ImageData), true, nil
	}

	return ImageData{}, false, nil
}

func (i *Img) downloadImage(url string) (ImageData, error) {
	var imageData ImageData
	res, err := http.Get(url)

	if err != nil {
		return imageData, err
	}

	defer res.Body.Close()

	src, err := imaging.Decode(res.Body)

	if err != nil {
		return imageData, err
	}

	imgType, err := extractImgType(res.Header)

	if err != nil {
		return imageData, err
	}

	return ImageData{
		ImgType: imgType,
		Header:  res.Header,
		image:   &src,
	}, nil
}

func extractImgType(header http.Header) (string, error) {
	content := header.Get("Content-Type")
	s := strings.Split(content, "/")

	if len(s) >= 2 {
		return s[1], nil
	} else {
		return "", errors.New("unable to determine image format")
	}
}
