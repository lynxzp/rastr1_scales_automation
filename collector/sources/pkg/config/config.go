package config

import (
	"collector/pkg/time"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	ScalesNums = 15
)

type Configurable struct {
	ListenPort int             `json:"listen_port"`
	ListenIP   string          `json:"listen_ip"`
	Users      map[string]User `json:"users"`
	Shifts     []Shift         `json:"shifts"`
}

var Cfg Configurable

type User struct {
	Password             string `json:"password"`
	AccessChangeFraction bool   `json:"accessChangeFraction"`
}

type Shift struct {
	Number int       `json:"number"`
	Start  time.Time `json:"start"`
	Finish time.Time `json:"finish"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Cfg.Users = make(map[string]User)

	f, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("WW cant' open config:", err)
		return
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("WW can't read log file:", err)
	}

	err = json.Unmarshal(bs, &Cfg)
	if err != nil {
		log.Fatalln("WW can't parse data in file", err)
	}
}
