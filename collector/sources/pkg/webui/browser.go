package webui

import (
	"collector/pkg/ucma"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	//pingPeriod = (pongWait * 9) / 10 // original time-out
	pingPeriod = 2 * time.Second // just in case network problem

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var (
	config Config
	scales *[ucma.ScalsesNums]ucma.Ucma
)

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	if containsDotDot(r.URL.Path) {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	http.ServeFile(w, r, "pkg/webui/www/"+r.URL.Path)
}

func ajax_update(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < ucma.ScalsesNums; i++ {
		scales[i].Ready = false

		dtype, ok := r.URL.Query()["dtype"+strconv.Itoa(i)]
		if !ok {
			continue
		}
		val, err := strconv.ParseUint(dtype[0], 16, 8)
		if err != nil {
			continue
		}
		scales[i].DType = uint16(val)

		ipaddr, ok := r.URL.Query()["ipaddr"+strconv.Itoa(i)]
		if !ok || len(ipaddr[0]) < 7 {
			continue
		}
		scales[i].IP = ipaddr[0]

		rs485addr, ok := r.URL.Query()["rs485addr"+strconv.Itoa(i)]
		if !ok || len(rs485addr[0]) < 1 {
			continue
		}
		val, err = strconv.ParseUint(rs485addr[0], 10, 8)
		if err != nil {
			continue
		}
		scales[i].Rs485addr = uint8(val)

		scales[i].Ready = true
	}

	// response
	log.Println(scales[0].Data)
	js, err := json.Marshal(scales)
	if err != nil {
		log.Println(err)
	}
	w.Write(js)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/ajax_update" {
		ajax_update(w, r)
		return
	}
	if r.URL.Path != "/" {
		serveFile(w, r)
		return
	}
	http.ServeFile(w, r, "pkg/webui/www/index.html")
}

func StartWeb(config_ Config, sc *[ucma.ScalsesNums]ucma.Ucma) {
	scales = sc
	config = config_
	http.HandleFunc("/", serveHome)
	log.Println(config.ListenIP)
	log.Println(config.ListenPort)
	log.Fatal(http.ListenAndServe(config.ListenIP+":"+config.ListenPort, nil))
}
