package settings

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"log"
)

type Settings struct {
	Ip      string
	Cname   string
	Port    string
	Address string
	Network string
}

func NewSettings() *Settings {
	return &Settings{"", "", "", "", "tcp"}
}
func (settings *Settings) LoadSettings() {
	file := os.Args[1]
	data, err := ioutil.ReadFile(file)
	err = json.Unmarshal(data, &settings)
	if err != nil {
		log.Fatalln(err)
	}
	settings.Address = settings.Ip + settings.Port
}
