package main

import (
	"collector/pkg/ucma"
	"collector/pkg/webui"
	"log"
	"sync/atomic"
	"time"
)

var Scales [ucma.ScalsesNums]ucma.Ucma

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go webui.StartWeb(cfg.webui, &Scales)

	time.Sleep(1 * time.Second)
	webui.OpenBrowser("http://127.0.0.1:8080")
	requestDelay := 1000*time.Millisecond
	for _, sc := range Scales {
		sc.Start(requestDelay)
	}
	for {
		log.Println(atomic.LoadInt32(&Scales[0].Data))
		<-time.After(requestDelay)
	}
}
