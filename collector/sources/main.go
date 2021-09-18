package main

import (
	"collector/pkg/ucma"
	"collector/pkg/webui"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	inChan := make(chan []byte)
	outChan := make(chan []byte)
	go webui.StartWeb(cfg.webui, inChan, outChan)

	time.Sleep(1*time.Second)
	webui.OpenBrowser("http://127.0.0.1:8080")
	u1 := ucma.Ucma{}
	u1.Start("192.168.1.12","502", 1*time.Millisecond)
	for {

	}
}

