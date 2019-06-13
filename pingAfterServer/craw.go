package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	Config = &config{
		HostList: viper.GetStringMapString(`agent`),
	}
}

var Config *config

type config struct {
	HostList map[string]string
}

type AgentJson struct {
	HtoH          string `json: "hostname"`
	ErrCount      string `json: "errCount"`
	Response_Time int64  `json: "response_Time"`
}

type ServerJson struct {
	Status    string       `json: "status"`
	AgentInfo []*AgentJson `json: "agentinfo"`
}

type Run struct {
	NodeStatus map[string]*ServerJson
}

func (this *Run) CrawInfo(serverjson *ServerJson, hostname, ip string, client *http.Client) {
	resp, err := client.Get("http://" + ip + "/info")
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		serverjson.Status = fmt.Sprintf("%s", err)
		this.NodeStatus[hostname] = serverjson
		return
	}
	jsonres, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		serverjson.Status = fmt.Sprintf("%s", err)
		this.NodeStatus[hostname] = serverjson
		return
	}

	if err := json.Unmarshal([]byte(jsonres), &serverjson.AgentInfo); err != nil {
		serverjson.Status = fmt.Sprintf("%s", err)
		this.NodeStatus[hostname] = serverjson
		return
	}

	serverjson.Status = "success"
	this.NodeStatus[hostname] = serverjson
}

func (this *Run) RunCrawInfo() {
	this.NodeStatus = make(map[string]*ServerJson)
	serverjson := ServerJson{}

	client := &http.Client{
		Timeout: time.Duration(1 * time.Second),
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			for hostname, ip := range Config.HostList {
				go this.CrawInfo(&serverjson, hostname, ip, client)
			}
		}
	}
}

func main() {
	run := &Run{}

	go run.RunCrawInfo()

	infohttp := func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(run.NodeStatus); err != nil {
			log.Println(err)
		}
	}
	http.HandleFunc("/info", infohttp)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
