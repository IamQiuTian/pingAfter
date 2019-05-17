package main

import (
	"log"

	"github.com/spf13/viper"
)

var Config *config

type config struct {
	HostList     map[string]string
	Hostname     string
	Alert_script string
	Execute      string
	Interval     int64
	Corrtime     int64
	Timer        int64
	listen         string
	To           string
}

func init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	Config = &config{
		HostList:     viper.GetStringMapString(`ip`),
		Hostname:     viper.GetString(`local.hostname`),
		Alert_script: viper.GetString(`alert.alert_script`),
		Execute:      viper.GetString(`alert.execute`),
		To:           viper.GetString(`alert.to`),
		Interval:     viper.GetInt64(`alert.interval`),
		Corrtime:     viper.GetInt64(`alert.corrtime`),
		Timer:        viper.GetInt64(`task.timer`),
		listen:         viper.GetString(`http.listen`),
	}
}
