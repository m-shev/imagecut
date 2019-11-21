package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
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
	excluded, _ := a.cache.Set(key, data, data.Size)

	for _, v := range excluded {
		err := os.Remove(v.(img.ImageData).Path)

		if err != nil {
			fmt.Println("remove error", err)
			//Todo log error
		}
	}
}

func (a *Api) flushCache() error {
	cache := a.cache.Flush()

	bytes, err := json.Marshal(cache)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(a.cachePath, bytes, 0644)

	return err
}

func (a *Api) removeCacheFile() error {
	_, err := os.Stat(a.cachePath)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	err = os.Remove(a.cachePath)

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
