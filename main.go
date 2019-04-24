package main

import (
    "log"
	"time"
)

func main() {
    log.Println("Start....")
	for {
		select {
		case <-time.After(time.Duration(Config.Timer) * time.Second):
			Run()
		}
	}
}
