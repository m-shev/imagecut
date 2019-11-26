package img

import (
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Img struct {
	imageFolder     string
	downloadTimeout time.Duration
}

var supportedImgTypes = []string{"jpeg", "png", "gif", "tiff", "bmp"}
var supportedImgTypesStr = strings.Join(supportedImgTypes, ", ")

type ImageData struct {
	ImgType    string
	Path       string
	Size       uint
	StatusCode int
	Header     http.Header
	src        io.Reader
}

type ImageSource struct {
	Url      string
	Headers  *http.Header
	FileName string
}

func NewImg(imageFolder string, downloadTimeout time.Duration) *Img {
	return &Img{
		imageFolder:     imageFolder,
		downloadTimeout: downloadTimeout,
	}
}

func (i *Img) CropByUrl(source ImageSource, width, height int) (ImageData, error) {
	data, err := i.downloadFile(source.Url, source.Headers)

	if err != nil {
		return data, err
	}

	if err = isImgTypeSupported(data.ImgType); err != nil {
		return data, err
	}

	im, err := imaging.Decode(data.src)

	if err != nil {
		return data, err
	}

	im = imaging.CropAnchor(im, width, height, imaging.Center)
	data.Path = fmt.Sprintf("%s/%s.%s", i.imageFolder, source.FileName, data.ImgType)

	err = imaging.Save(im, data.Path)

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

func isImgTypeSupported(imgType string) error {
	for _, v := range supportedImgTypes {
		if v == imgType {
			return nil
		}
	}

	return makeUnsupportedImgFormatError(imgType)
}

func makeUnsupportedImgFormatError(imgType string) error {
	return fmt.Errorf(
		"unsupported image format: \"%s\", avalibale: %s",
		imgType,
		supportedImgTypesStr,
	)
}
