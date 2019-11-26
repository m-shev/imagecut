package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imagecut/internal/img"
	"imagecut/internal/lru"
	"io/ioutil"
	"os"
)

type CacheItem struct {
	Key   string
	Size  uint
	Value img.ImageData
}

func (a *Api) getFromCache(key string, ctx *gin.Context) (img.ImageData, bool) {
	a.Mutex.Lock()
	v, err := a.cache.Get(key)
	a.Mutex.Unlock()

	a.logOnErr(ctx, err)

	if v != nil {
		return v.(img.ImageData), true
	}

	return img.ImageData{}, false
}

func (a *Api) setToCache(key string, data img.ImageData, ctx *gin.Context) {

	a.Mutex.Lock()
	excluded, err := a.cache.Set(key, data, data.Size)
	a.Mutex.Unlock()

	a.logOnErr(ctx, err)

	for _, v := range excluded {
		err := os.Remove(v.(img.ImageData).Path)
		a.logOnErr(ctx, err)
	}
}

func (a *Api) FlushCache() error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	cache := a.cache.Flush()

	bytes, err := json.Marshal(cache)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(a.cachePath, bytes, 0644)

	return err
}

func (a *Api) restoreCache() error {
	queue := make([]CacheItem, 0)

	bytes, err := ioutil.ReadFile(a.cachePath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &queue)

	if err != nil {
		return err
	}

	lruQueue := make([]*lru.CacheItem, len(queue))

	for index, v := range queue {
		lruQueue[index] = &lru.CacheItem{
			Value: v.Value,
			Key:   v.Key,
			Size:  v.Size,
		}
	}

	a.cache.RestoreData(lruQueue)

	return nil
}

func hasher(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
