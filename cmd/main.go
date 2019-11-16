package main

import (
	"github.com/gin-gonic/gin"
	"imagecut/api"
	"imagecut/config"
	"log"
	"net/http"
)

func main() {
	config.AddConfigPath("../config")
	conf := config.GetConfig()
	api := api.NewApi()
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
	handler.GET("/crop/:height/:width/*url", api.Crop)

	server := &http.Server{
		Addr:    conf.Http.Addr,
		Handler: handler,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
