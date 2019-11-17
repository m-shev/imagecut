package lru

import (
	"errors"
	"fmt"
	"imagecut/internal/linkedlist"
)

type Lru struct {
	size       uint
	maxSize    uint
	list linkedlist.List
	cacheMap map[string]*linkedlist.Item
}

type CacheItem struct {
	value interface{}
	key string
	size uint
}

func NewLru(maxSize uint) *Lru {
	return &Lru{
		size:     0,
		maxSize:  maxSize,
		list:     linkedlist.List{},
		cacheMap: make(map[string]*linkedlist.Item),
	}
}

func (lru *Lru) Set(key string, value interface{}, size uint) ([]interface{}, error) {
	_, ok := lru.cacheMap[key]

	if ok {
		return nil, errors.New(fmt.Sprintf("value with key %s has already been added to the cache", key))
	} else {
		item := lru.list.PushFront(&CacheItem{
			value: value,
			size:  size,
			key: key,
		})

		lru.cacheMap[key] = item
		lru.size += size
		exclusion, err := lru.cleanCache()
		return exclusion, err
	}
}

func (lru *Lru) Get(key string) (interface{}, error) {
	if item, ok := lru.cacheMap[key]; ok {
		err := item.Remove()

		if err != nil {
			return nil, err
		}

		cacheItem := item.Value().(*CacheItem)
		item := lru.list.PushFront(cacheItem)
		lru.cacheMap[key] = item

		return cacheItem.value, nil
	}

	return nil, nil
}

func (lru *Lru) cleanCache() ([]interface{}, error) {
	excludedItems := make([]interface{}, 0)

	if lru.size > lru.maxSize {
		for lru.size > lru.maxSize {
			item := lru.list.Last()
			err := item.Remove()

			if err != nil {
				return excludedItems, err
			}

			cacheItem := item.Value().(*CacheItem)

			delete(lru.cacheMap, cacheItem.key)

			lru.size -= cacheItem.size
			excludedItems = append(excludedItems, cacheItem.value)
		}
	}

	return excludedItems, nil
}

