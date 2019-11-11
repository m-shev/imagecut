package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"imagecut/config"
	"log"
	"net/http"
)

func main()  {
	config.AddConfigPath("../config")
	conf := config.GetConfig()

	gin.DefaultWriter = &lumberjack.Logger{
		Filename:   "foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	gin.Recovery()
	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(gin.Logger())
	handler.Use(testMid())


	handler.GET("/status", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
		fmt.Println("handled")
	})




	err := handler.Run(conf.Http.Addr)
	//server := &http.Server{
	//	Addr:              conf.Http.Addr,
	//	Handler:           handler,
	//}
	//
	//err = server.ListenAndServe()
	//
	if err != nil {
		log.Fatal(err)
	}
}

func testMid() func (ctx *gin.Context){
	return func(ctx *gin.Context) {
		fmt.Println("after")
		ctx.Next()
	}
}