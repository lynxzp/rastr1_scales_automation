package main

import (
	"collector/pkg/ucma"
	"collector/pkg/webui"
	"log"
	"time"
)

var Scales [ucma.ScalsesNums]ucma.Ucma

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go webui.StartWeb(cfg.webui, &Scales)

	time.Sleep(1 * time.Second)
	webui.OpenBrowser("http://127.0.0.1:8080")
	requestDelay := 1000*time.Millisecond
	for i := range Scales {
		Scales[i].Start(requestDelay)
	}
	for {
		<-time.After(requestDelay)
	}
}
