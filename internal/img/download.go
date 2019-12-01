package img

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (i *Img) downloadFile(url string, header *http.Header) (ImageData, error) {
	var imageData ImageData

	ctx, cancel := context.WithTimeout(context.Background(), i.downloadTimeout*time.Second)
	defer cancel()

	res, err := i.makeReq(ctx, url, header)

	if err != nil {
		return imageData, err
	}

	if res.StatusCode != 200 && res.StatusCode != 304 {
		imageData.StatusCode = res.StatusCode
		if res.StatusCode == 404 {
			err = fmt.Errorf("image not found")
		} else {
			err = fmt.Errorf("cannot download image")
		}
		return imageData, err
	}

	defer res.Body.Close()

	imgType, err := extractImgType(res.Header)

	if err != nil {
		return imageData, err
	}

	src := bytes.NewBuffer(make([]byte, 0, getContentLength(res.Header)))
	_, err = io.Copy(src, res.Body)

	if err != nil {
		return imageData, err
	}

	return ImageData{
		ImgType:    imgType,
		Header:     res.Header,
		StatusCode: res.StatusCode,
		src:        src,
	}, nil
}

func (i *Img) makeReq(ctx context.Context, url string, header *http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	addHeaders(req, header)

	client := http.DefaultClient
	res, err := client.Do(req)

	return res, err
}

func addHeaders(req *http.Request, headers *http.Header) {
	for k, v := range *headers {
		for _, h := range v {
			if k != "If-Modified-Since" && k != "If-None-Match" {
				req.Header.Set(k, h)
			}
		}
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
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

func getContentLength(header http.Header) int {
	l := header.Get("Content-Length")
	contentLength, err := strconv.Atoi(l)

	if err != nil {
		return 0
	}

	return contentLength
}
