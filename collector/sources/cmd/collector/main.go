package main

import (
	"collector/pkg/config"
	"collector/pkg/ucma"
	"collector/pkg/webui"
	"log"
	"strconv"
	"time"
)

var Scales [config.ScalesNums]ucma.Ucma

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go webui.StartWeb(&Scales)

	time.Sleep(1 * time.Second)
	webui.OpenBrowser("http://" + config.ListenIP + ":" + strconv.Itoa(config.ListenPort))
	requestDelay := 1000 * time.Millisecond
	for i := range Scales {
		Scales[i].Id = int8(i)
		Scales[i].Start(requestDelay)
	}
	for {
		<-time.After(requestDelay)
	}
}
