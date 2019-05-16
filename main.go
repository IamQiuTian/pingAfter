package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	log.Println("Start....")
	go func() {
		r := gin.Default()
		r.GET("/info", Jsonapi)
		r.Run(fmt.Sprintf(":%s", Config.Port))
	}()
	for {
		select {
		case <-time.After(time.Duration(Config.Timer) * time.Second):
			Run()
		}
	}
}
