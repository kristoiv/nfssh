package main

import (
	"encoding/json"
	"io/ioutil"
)

var Config *config

func init() {
	Config = &config{}

	rawJson, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(rawJson, Config); err != nil {
		panic(err)
	}
}

type config struct {
	Host                 string
	Username             string
	LocalNfsdPort        int
	RemoteNfsdPort       int
	LocalMountdPort      int
	RemoteMountdPort     int
	LocalVolumeName      string
	RemoteMountDirectory string
}
