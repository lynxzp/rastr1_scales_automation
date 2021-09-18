package main

import "collector/pkg/webui"

type config struct {
	webui webui.Config
}

var cfg config

func init() {
	cfg.webui.ListenIP = "0.0.0.0"
	cfg.webui.ListenPort = "8080"

}