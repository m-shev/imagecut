package lru

import (
	"encoding/json"
	"errors"
	"fmt"
	"imagecut/internal/linkedlist"
	"io/ioutil"
)

type Lru struct {
	path     string `json:"_"`
	size     uint
	maxSize  uint `json:"_"`
	list     linkedlist.List
	cacheMap map[string]*linkedlist.Item
}


type CacheData struct {
	size     uint
	maxSize  uint
	list     linkedlist.List
	cacheMap map[string]*linkedlist.Item
}

type CacheItem struct {
	value interface{}
	key   string
	size  uint
}

func NewLru(maxSize uint, path string) *Lru {
	return &Lru{
		path: path,
		size:     0,
		maxSize:  maxSize,
		list:     linkedlist.List{},
		cacheMap: make(map[string]*linkedlist.Item),
	}
}

func (l *Lru) Set(key string, value interface{}, size uint) ([]interface{}, error) {
	_, ok := l.cacheMap[key]

	if ok {
		return nil, errors.New(fmt.Sprintf("value with key %s has already been added to the cache", key))
	} else {
		item := l.list.PushFront(&CacheItem{
			value: value,
			size:  size,
			key:   key,
		})

		l.cacheMap[key] = item
		l.size += size
		exclusion, err := l.cleanCache()
		return exclusion, err
	}
}

func (l *Lru) Get(key string) (interface{}, error) {
	if item, ok := l.cacheMap[key]; ok {
		err := item.Remove()

		if err != nil {
			return nil, err
		}

		cacheItem := item.Value().(*CacheItem)
		item := l.list.PushFront(cacheItem)
		l.cacheMap[key] = item

		return cacheItem.value, nil
	}

	return nil, nil
}

//func (l *Lru) readData {
//
//}

func (l *Lru) Flush() error {
	str, err := json.Marshal(l)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(l.path, str, 0644)

	return err
}

func (l *Lru) cleanCache() ([]interface{}, error) {
	excludedItems := make([]interface{}, 0)

	if l.size > l.maxSize {
		for l.size > l.maxSize {
			item := l.list.Last()
			err := item.Remove()

			if err != nil {
				return excludedItems, err
			}

			cacheItem := item.Value().(*CacheItem)

			delete(l.cacheMap, cacheItem.key)

			l.size -= cacheItem.size
			excludedItems = append(excludedItems, cacheItem.value)
		}
	}

	return excludedItems, nil
}
