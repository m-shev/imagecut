package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"imagecut/api"
	"imagecut/config"
	"imagecut/internal/lru"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.AddConfigPath("../config")
	conf := config.GetConfig()
	cache, err := lru.NewLru(conf.CacheSize, conf.CachePath)

	if err !=nil {
		log.Fatal(err)
	}

	api := api.NewApi(cache, conf.ImageFolder)
	//gin.DefaultWriter = &lumberjack.Logger{
	//	Filename:   "foo.log",
	//	MaxSize:    500, // megabytes
	//	MaxBackups: 3,
	//	MaxAge:     28, // days
	//}


	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(gin.Logger())

	handler.GET("/status", api.Status)
	handler.GET("/crop/:height/:width/", api.Crop)

	server := &http.Server{
		Addr:    conf.Http.Addr,
		Handler: handler,
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	graceful(server, 5 * time.Second, func() {
		err := cache.Flush()
		log.Println(err)
	})
}

func graceful(hs *http.Server, timeout time.Duration, callback func()) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	callback()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}