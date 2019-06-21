package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/json-iterator/go"
	"github.com/spf13/viper"
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

type AgentJsonInfo struct {
	HtoH          string `json: "htoh"`
	ErrCount      string `json: "errCount"`
	Response_Time int64  `json: "response_Time"`
}

type ServerJson struct {
	Status    string           `json: "status"`
	AgentInfo []*AgentJsonInfo `json: "agentinfo"`
}

// {"hk01 TO google":{"errCount":0,"response_Time":1}}
// {"hk02":{"Status":"success","AgentInfo":{"HTOH":"","JsonInfo":null}}}

type Run struct {
	NodeStatus map[string]ServerJson
}

var mux = &sync.Mutex{}
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (this *Run) CrawInfo(hostname, ip string) {
	mux.Lock()
	defer mux.Unlock()

	var err error
	var resp *http.Response
	var jsonres []byte

	serverJson := &ServerJson{}

	client := &http.Client{
		Timeout: time.Duration(1 * time.Second),
	}

	if resp, err = client.Get("http://" + ip + "/info"); err != nil {
		goto ERR
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if jsonres, err = ioutil.ReadAll(resp.Body); err != nil {
		goto ERR
	}

	if err := json.Unmarshal([]byte(jsonres), &serverJson.AgentInfo); err != nil {
		goto ERR
	}

	serverJson.Status = "success"
	this.NodeStatus[hostname] = *serverJson

ERR:
	serverJson.Status = fmt.Sprintf("%s", err)
	fmt.Println(err)
	fmt.Println(serverJson.AgentInfo)
	this.NodeStatus[hostname] = *serverJson
	return
}

func (this *Run) RunCrawInfo() {
	this.NodeStatus = make(map[string]ServerJson)

	for {
		select {
		case <-time.After(10 * time.Second):
			for hostname, ip := range Config.HostList {
				go this.CrawInfo(hostname, ip)
			}
		}
	}
}

func main() {
	run := &Run{}

	go run.RunCrawInfo()

	infohttp := func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&run.NodeStatus); err != nil {
			log.Println(err)
		}
	}
	http.HandleFunc("/info", infohttp)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
