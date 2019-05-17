package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	log.Println("Start....")
	go func() {
		r := gin.Default()
		r.GET("/info", Jsonapi)
		r.Run(Config.listen)
	}()
	for {
		select {
		case <-time.After(time.Duration(Config.Timer) * time.Second):
			Run()
		}
	}
}
