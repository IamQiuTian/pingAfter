package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os/exec"

	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var Info = make(map[string]map[string]int64)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Task(host string, w *sync.WaitGroup) {
	defer w.Done()

	PingValue := map[string]int64{"errCount": 0, "response_Time": 0}

	raddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return
	}
	for i := 1; i < int(Config.Interval)+1; i++ {
		response_time, err := SendICMPRequest(GetICMP(uint16(i)), raddr)
		if err != nil {
			log.Println(err)
			PingValue["errCount"] = PingValue["errCount"] + 1
			continue
		}
		PingValue["response_Time"] += response_time

		time.Sleep(1 * time.Second)
	}

	PingValue["response_Time"] = PingValue["response_Time"] / Config.Interval
	message := fmt.Sprintf("Number of errors: %s,  Response time: %sms", strconv.Itoa(int(PingValue["errCount"])), strconv.Itoa(int(PingValue["response_Time"])))

	if PingValue["errCount"] >= Config.Interval || PingValue["response_Time"] >= Config.Corrtime {
		title := fmt.Sprintf("%s to %s(%s) ping error", Config.Hostname, Config.HostList[host], host)
		Afert(title, message)
	}


	Info[fmt.Sprintf("%s to %s", Config.Hostname, host)] = PingValue
	log.Printf("to %s : %s\n", host, message)
}

func Afert(title, message string) {
	cmd := exec.Command(Config.Execute, Config.Alert_script, Config.To, title, message)
	defer cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
}

func Jsonapi(c *gin.Context)  {
	c.JSON(200, &Info)
}

func Run() {
	for ip, _ := range Config.HostList {
		wg.Add(1)
		go Task(ip, &wg)
	}
	wg.Wait()
}
