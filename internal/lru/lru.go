package lru

import (
	"encoding/json"
	"errors"
	"fmt"
	"imagecut/internal/linkedlist"
	"io/ioutil"
)

type Lru struct {
	path     string
	size     uint
	maxSize  uint
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
	Value interface{}
	Key   string
	Size  uint
}

func NewLru(maxSize uint, path string) (*Lru, error) {
	lru := &Lru{
		path: path,
		size:     0,
		maxSize:  maxSize,
		list:     linkedlist.List{},
		cacheMap: make(map[string]*linkedlist.Item),
	}

	err := lru.restoreData()

	if err != nil {
		return lru, err
	}

	return lru, nil
}

func (l *Lru) Set(key string, value interface{}, size uint) ([]interface{}, error) {
	_, ok := l.cacheMap[key]

	if ok {
		return nil, errors.New(fmt.Sprintf("Value with Key %s has already been added to the cache", key))
	} else {
		item := l.list.PushFront(&CacheItem{
			Value: value,
			Size:  size,
			Key:   key,
		})

		l.cacheMap[key] = item
		l.size += size
		excluded, err := l.cleanCache()
		return excluded, err
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

		return cacheItem.Value, nil
	}

	return nil, nil
}

func (l *Lru) restoreData() error {
	bytes, err := ioutil.ReadFile(l.path)

	if err != nil {
		return err
	}

	var queue []CacheItem

	err = json.Unmarshal(bytes, &queue)

	if err !=nil {
		return err
	}

	for _, cacheItem := range queue {
		item := l.list.PushBack(&cacheItem)
		l.cacheMap[cacheItem.Key] = item
		l.size += cacheItem.Size
	}

	return nil
}

func (l *Lru) Flush() error {
	queue := make([]CacheItem, 0)

	item := l.list.First()
		for item != nil {
			queue = append(queue, *item.Value().(*CacheItem))
			item = item.Next()
	}

	bytes, err := json.Marshal(queue)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(l.path, bytes, 0644)

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

			delete(l.cacheMap, cacheItem.Key)

			l.size -= cacheItem.Size
			excludedItems = append(excludedItems, cacheItem.Value)
		}
	}

	return excludedItems, nil
}
