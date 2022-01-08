package webui

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type passwords struct {
	password string
	//salt     string
	//hash     string
}

var users map[string]passwords

func init() {
	users = make(map[string]passwords)

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var salt string
	for i := 0; i < 16; i++ {
		salt += string(letterRunes[rand.Intn(len(letterRunes))])
	}
	login := "vasya"
	password := "password"
	users[login] = passwords{password}
}

func sendLoginForm(w http.ResponseWriter, r *http.Request, params string) {
	http.ServeFile(w, r, "pkg/webui/www/login.html")
}

func loginH(w http.ResponseWriter, r *http.Request) {
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

	if val, ok := users[login]; ok && (val.password == password) {
		cookie1 := &http.Cookie{
			Name:  "login",
			Value: login,
		}
		cookie2 := &http.Cookie{
			Name:  "password",
			Value: password,
		}
		http.SetCookie(w, cookie1)
		http.SetCookie(w, cookie2)
		serveMain(w, r)
	}
	sendLoginForm(w, r, "wrong=password")
	return
}

func loggined(r *http.Request) bool {
	log.Println(r.Cookies())
	var login, password string
	for _, c := range r.Cookies() {
		if c.Name == "login" {
			login = c.Value
		}
		if c.Name == "password" {
			password = c.Value
		}
	}
	log.Println(login, password)
	if val, ok := users[login]; ok && (val.password == password) {
		log.Println("loggined")
		return true
	}
	log.Println("access denied")
	return false
}
