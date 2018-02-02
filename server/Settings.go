package server

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"log"
	"project/settings"
)

func LoadSettings() settings.Settings {
	var s settings.Settings
	s.Network = "tcp"
	s.File = os.Args[1]
	data, err := ioutil.ReadFile(s.File)
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatalln(err)
	}
	s.Address = s.Ip + s.Port
	return s
}
