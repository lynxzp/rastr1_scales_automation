package webui

import (
	"collector/pkg/config"
	"collector/pkg/reports"
	"collector/pkg/shift"
	"collector/pkg/store"
	"collector/pkg/ucma"
	"encoding/json"
	"fmt"
	"log"
	"net"
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

func ajaxUpdate(w http.ResponseWriter) {
	js, err := json.Marshal(scales)
	if err != nil {
		log.Println(err)
	}
	_, _ = w.Write(js)
}

func ajaxSaveFraction(w http.ResponseWriter, r *http.Request) {
	if allowed, login := isAccessChangeFraction(r); allowed == false {
		log.Println(login, "tried to change fraction")
		return
	}
	idStr, ok := r.URL.Query()["id"]
	if !ok {
		return
	}
	id, err := strconv.ParseUint(idStr[0], 10, 8)
	if err != nil {
		return
	}

	fractionStr, ok := r.URL.Query()["fraction"]
	if !ok || len(fractionStr[0]) < 1 {
		return
	}

	store.SaveEvent(int(id), int(scales[id].DataAccumValue), "start", shift.GetCurrentShift(int8(id)), fractionStr[0])
	store.SaveScaleFraction(int(id), fractionStr[0])
	reloadScales()

	_, _ = w.Write([]byte("ok"))
}

func ajaxSaveScale(w http.ResponseWriter, r *http.Request) {
	allowed, login := isAccessChangeFraction(r)
	if allowed == false {
		log.Println(login, "tried to change fraction")
		return
	}
	if !isLocalhost(r) {
		log.Println(login, "tried to change scales")
		return
	}
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

	fractionStr, ok := r.URL.Query()["fraction"]
	if !ok || len(fractionStr[0]) < 1 {
		return
	}

	log.Println("new scales config:", id, ipaddr[0], rs485addr, fractionStr[0])
	store.SaveScale(int(id), int(dataPerfAddr), ipaddr[0], int(rs485addr), fractionStr[0])
	reloadScales()
	scales[id].Requests = 0
	scales[id].Responses = 0

	_, _ = w.Write([]byte("ok"))
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

	_, _ = w.Write([]byte("ok"))
}

func exportCSV(w http.ResponseWriter, r *http.Request) {
	sepParam, ok := r.URL.Query()["separator"]
	separator := ","
	if ok && sepParam[0] == ";" {
		separator = ";"
	}
	c := store.ExportData(separator)
	headers := w.Header()
	headers["Content-Type"] = []string{"text/csv"}
	for str := range c {
		_, _ = w.Write([]byte(str))
	}
}

type ReportParams []struct {
	Start        string         `json:"start"`
	End          string         `json:"end"`
	Shift        int            `json:"shift"`
	Column       string         `json:"column"`
	Accumulation map[string]int `json:"accumulation"`
}

func reportH(w http.ResponseWriter, r *http.Request) {
	layout := "02.01.2006 15:04:05"

	paramsStr, ok := r.URL.Query()["params"]
	if !ok {
		return
	}
	var params ReportParams
	err := json.Unmarshal([]byte(paramsStr[0]), &params)
	if err != nil {
		log.Println(err)
	}

	for i := range params {
		startTime, err := time.Parse(layout, params[i].Start)
		if err != nil {
			log.Println(err.Error())
			return
		}
		endTime, err := time.Parse(layout, params[i].End)
		if err != nil {
			log.Println(err.Error())
			return
		}
		params[i].Accumulation = reports.Count(startTime, endTime, params[i].Shift)
	}

	resp, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
	}

	_, _ = w.Write(resp)
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

func serve(w http.ResponseWriter, r *http.Request) {
	// Delete any ETag headers that may have been set
	for _, v := range etagHeaders {
		if r.Header.Get(v) != "" {
			r.Header.Del(v)
		}
	}
	// Set our NoCache headers
	for k, v := range noCacheHeaders {
		w.Header().Set(k, v)
	}

	if !loggined(r) {
		loginH(w, r)
		return
	}
	if r.URL.Path == "/ajax_update" {
		ajaxUpdate(w)
		return
	}
	if r.URL.Path == "/save_fraction" {
		ajaxSaveFraction(w, r)
		return
	}
	if r.URL.Path == "/save_scale" {
		ajaxSaveScale(w, r)
		return
	}
	if r.URL.Path == "/clear" {
		ajaxClear(w, r)
		return
	}
	if r.URL.Path == "/export" {
		exportCSV(w, r)
		return
	}
	if r.URL.Path == "/report" {
		reportH(w, r)
		return
	}
	if r.URL.Path == "/login" {
		loginH(w, r)
		http.Redirect(w, r, "/", 200)
		return
	}
	if r.URL.Path != "/" {
		serveFile(w, r)
		return
	}
	serveMain(w, r)
}

func serveMain(w http.ResponseWriter, r *http.Request) {
	if isLocalhost(r) {
		http.ServeFile(w, r, "pkg/webui/www/admin.html")
		return
	}
	http.ServeFile(w, r, "pkg/webui/www/index.html")
}

func StartWeb(sc *[config.ScalesNums]ucma.Ucma) {
	scales = sc
	reloadScales()
	http.HandleFunc("/", serve)
	log.Fatal(http.ListenAndServe(config.Cfg.ListenIP+":"+strconv.Itoa(config.Cfg.ListenPort), nil))
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
		scales[i].Fraction = s[i].Fraction
	}
}

func isLocalhost(r *http.Request) bool {
	if ip, err := getIP(r); err == nil && ip == "127.0.0.1" {
		return true
	}
	return false
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}
