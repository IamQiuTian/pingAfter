package main

import (
	"time"
)

func main() {
	for {
		select {
		case <-time.After(time.Duration(Config.Timer) * time.Second):
			Run()
		}
	}
}
