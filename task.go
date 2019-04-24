package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Task(host string, w *sync.WaitGroup) {
	PingValue := map[string]int64{"errCount": 0, "response_Time": 0}

	raddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Println(err)
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
		title := fmt.Sprintf("%s to %s ping error", Config.Hostname, host)
		Afert(title, message)
	}
	log.Printf("to %s : %s\n", host, message)
	defer w.Done()
}

func Afert(title, message string) {
	cmd := exec.Command(Config.Execute, Config.Alert_script, Config.To, title, message)
	defer cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
}

func Run() {
	for _, ip := range Config.IpList {
		wg.Add(1)
		go Task(ip, &wg)
	}
	wg.Wait()
}
