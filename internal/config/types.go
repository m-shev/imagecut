package config

import "time"

type Config struct {
	Http
	Img
	Cache
	Logging
	env    string
	isRead bool
}

type Http struct {
	Addr string
}

type Img struct {
	ImageFolder     string
	DownloadTimeout time.Duration
}

type Cache struct {
	CacheSize uint
	CachePath string
}

type Logging struct {
	ErrorLog  LogParams
	AccessLog LogParams
}

type LogParams struct {
	FileName   string
	MaxBackups int
	MaxAge     int
}
