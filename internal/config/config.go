package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
)

type Config struct {
	Ignored []string `json:"ignored"`
	Command string   `json:"command"`
}

func GetConfig(loc string) Config {
	content, err := ioutil.ReadFile(path.Join(loc, "/.stalkerrc.json"))
	if err != nil {
		return Config{}
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}
