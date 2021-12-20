package webui

import (
	"collector/pkg/config"
	"collector/pkg/store"
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
	scales *[config.ScalesNums]ucma.Ucma
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

func ajaxUpdate(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(scales)
	if err != nil {
		log.Println(err)
	}
	w.Write(js)
}

func ajaxSave(w http.ResponseWriter, r *http.Request) {
	idStr, ok := r.URL.Query()["id"]
	if !ok {
		return
	}
	id, err := strconv.ParseUint(idStr[0], 10, 8)
	if err != nil {
		return
	}

	dataPerfAddrStr, ok := r.URL.Query()["dtype"]
	if !ok {
		return
	}
	dataPerfAddr, err := strconv.ParseUint(dataPerfAddrStr[0], 16, 8)
	if err != nil {
		return
	}

	ipaddr, ok := r.URL.Query()["ipaddr"]
	if !ok || len(ipaddr[0]) < 7 {
		return
	}

	rs485addrStr, ok := r.URL.Query()["rs485addr"]
	if !ok || len(rs485addrStr[0]) < 1 {
		return
	}
	rs485addr, err := strconv.ParseUint(rs485addrStr[0], 10, 8)
	if err != nil {
		return
	}

	store.SaveScale(int(id), int(dataPerfAddr), ipaddr[0], int(rs485addr))
	reloadScales()
	scales[id].Requests = 0
	scales[id].Responses = 0

	w.Write([]byte("ok"))
}

func ajaxClear(w http.ResponseWriter, r *http.Request) {
	idStr, ok := r.URL.Query()["id"]
	if !ok {
		return
	}
	id, err := strconv.ParseUint(idStr[0], 10, 8)
	if err != nil {
		return
	}

	store.ClearScale(int(id))
	scales[id].Requests = 0
	scales[id].Responses = 0
	scales[id].DataPerfValue = 0
	scales[id].DataAccumValue = 0
	reloadScales()

	w.Write([]byte("ok"))
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	c := store.ExportData()
	headers := w.Header()
	headers["Content-Type"] = []string{"text/csv"}
	for str := range c {
		w.Write([]byte(str))
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/ajax_update" {
		ajaxUpdate(w, r)
		return
	}
	if r.URL.Path == "/save" {
		ajaxSave(w, r)
		return
	}
	if r.URL.Path == "/clear" {
		ajaxClear(w, r)
		return
	}
	if r.URL.Path == "/export" {
		exportHandler(w, r)
		return
	}
	if r.URL.Path != "/" {
		serveFile(w, r)
		return
	}
	http.ServeFile(w, r, "pkg/webui/www/index.html")
}

func StartWeb(sc *[config.ScalesNums]ucma.Ucma) {
	scales = sc
	reloadScales()
	http.HandleFunc("/", serveHome)
	log.Fatal(http.ListenAndServe(config.ListenIP+":"+strconv.Itoa(config.ListenPort), nil))
}

func reloadScales() {
	s, err := store.ReadScales()
	if err != nil {
		log.Println("EE Can't load scales config: ", err)
		return
	}
	if len(s) != config.ScalesNums {
		log.Println("WW config.ScalesNums=", config.ScalesNums, " scales in db=", len(s))
		return
	}
	for i := range s {
		scales[i].DataPerfAddr = uint16(s[i].DataPerfAddr)
		scales[i].Rs485addr = uint8(s[i].Rs485addr)
		scales[i].IP = s[i].Ip
	}
}
