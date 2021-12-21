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
	ListenPort int     `json:"listen_port"`
	ListenIP   string  `json:"listen_ip"`
	Users      []User  `json:"users"`
	Shifts     []Shift `json:"shifts"`
}

var Cfg Configurable

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Rights   string `json:"rights"`
}

type Shift struct {
	Number int       `json:"number"`
	Start  time.Time `json:"start"`
	Finish time.Time `json:"finish"`
}

func init() {
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

	for i := range Cfg.Users {
		switch Cfg.Users[i].Rights {
		case "read_only":
		case "change_fraction":
		default:
			log.Println("EE wrong user rights:", Cfg.Users[i])
			Cfg.Users[i].Rights = "read_only"
		}
	}
}
