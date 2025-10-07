package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/tidwall/jsonc"
)

type Account struct {
	Email string              `json:"email"`
	Key   string              `json:"key"`
	Zones map[string][]string `json:"zones"`
}

type UnifiSettings struct {
	Enable  bool   `json:"enable"`
	Gateway string `json:"gateway"`
	Token   string `json:"token"`
	ListID  string `json:"listID"`
}

type Settings struct {
	Interval uint          `json:"interval"`
	IPV6     bool          `json:"ipv6"`
	Unifi    UnifiSettings `json:"unifi"`
	Accounts []Account     `json:"accounts"`
}

var (
	settings Settings
)

func loadSettings() {
	data, err := os.ReadFile("settings.jsonc")
	if err == nil {
		data = jsonc.ToJSONInPlace(data)
	} else {
		data, err = os.ReadFile("settings.json")
		if err != nil {
			log.Panicln("failed to read settings.json(c)")
		}
	}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		log.Panicln(err)
	}
}
