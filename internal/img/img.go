package img

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"net/http"
	"os"
	"strings"
)

type Img struct {
	imageFolder string
}

type ImageData struct {
	ImgType string
	Path    string
	Size    uint
	Header  http.Header
	image   *image.Image
}

func NewImg(imageFolder string) *Img {
	return &Img{imageFolder: imageFolder}
}

func (i *Img) CropByUrl(url, fileName string, width, height int) (ImageData, error) {
	data, err := i.downloadImage(url)

	if err != nil {
		return data, err
	}

	cropped := imaging.CropAnchor(*data.image, width, height, imaging.Center)
	data.Path = fmt.Sprintf("%s/%s.%s", i.imageFolder, fileName, data.ImgType)

	err = imaging.Save(cropped, data.Path)

	if err != nil {
		return data, err
	}

	err = setFileSize(&data)

	return data, nil
}

func setFileSize(data *ImageData) error {
	stat, err := os.Stat(data.Path)

	if err != nil {
		return err
	}

	data.Size = uint(stat.Size())

	return nil
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
