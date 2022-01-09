package webui

import (
	"collector/pkg/config"
	"fmt"
	"net/http"
	"strconv"
)

func sendLoginForm(w http.ResponseWriter, r *http.Request, params string) {
	http.ServeFile(w, r, "pkg/webui/www/login.html")
}

func loginH(w http.ResponseWriter, r *http.Request) {
	if loggined(r) {
		redirectToMain(w)
		return
	}
	if r.Method != "POST" {
		sendLoginForm(w, r, "wrong=method")
		return
	}
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseFrom() err: %v", err)
		return
	}
	login := r.FormValue("login")
	password := r.FormValue("password")

	if val, ok := config.Cfg.Users[login]; ok && (val.Password == password) {
		cookie1 := &http.Cookie{
			Name:  "login",
			Value: login,
		}
		cookie2 := &http.Cookie{
			Name:  "password",
			Value: password,
		}
		cookie3 := &http.Cookie{
			Name:  "AccessChangeFraction",
			Value: strconv.FormatBool(config.Cfg.Users[login].AccessChangeFraction),
		}
		http.SetCookie(w, cookie1)
		http.SetCookie(w, cookie2)
		http.SetCookie(w, cookie3)

		redirectToMain(w)
		return
	}
	sendLoginForm(w, r, "wrong=password")
	return
}

func loggined(r *http.Request) bool {
	var login, password string
	for _, c := range r.Cookies() {
		if c.Name == "login" {
			login = c.Value
		}
		if c.Name == "password" {
			password = c.Value
		}
	}
	if val, ok := config.Cfg.Users[login]; ok && (val.Password == password) {
		return true
	}
	return false
}

func isAccessChangeFraction(r *http.Request) bool {
	var login, password string
	for _, c := range r.Cookies() {
		if c.Name == "login" {
			login = c.Value
		}
		if c.Name == "password" {
			password = c.Value
		}
	}
	if val, ok := config.Cfg.Users[login]; !ok || (val.Password != password) {
		return false
	}

	return config.Cfg.Users[login].AccessChangeFraction
}

func redirectToMain(w http.ResponseWriter) {
	_, _ = w.Write([]byte(`<html>
    <head>
        <meta http-equiv="refresh" content="1;url=/" />
    </head>
    <body>
        <h1>Вход успешен. Загрузка...</h1>
    </body>
</html>
`))
}
