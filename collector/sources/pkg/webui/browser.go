package webui

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/grafov/bcast"
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
)


// Transmit input from websocket to specified program.
func pumpStdin(ws *websocket.Conn, cmdChan chan []byte) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		cmdChan <- message
	}
}

// Send data to websocket.
func sendWebsocket(ws *websocket.Conn, b []byte) error {

	err := ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		ws.Close()
		return err
	}
	err = ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		ws.Close()
		return err
	}
	return nil
}

// Transmit data from stdout specified program to websocket.
func pumpStdout(ws *websocket.Conn, bCastMember *bcast.Member, done chan struct{}) {
	defer close(done)
	defer ws.Close()

	for {
		buf := bCastMember.Recv()
		err := sendWebsocket(ws, buf.([]byte))
		if err == io.EOF {
			fmt.Println("\nEOF in PTY")
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
}

// Pinging websocket (for supporting connection)
func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				if err.Error() =="websocket: close sent" {
					// workaround, auto close ping function did not work, so will exit in this way
					return
				}
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

var upgrader = websocket.Upgrader{}

// Http websocket handler.
func serveWs(w http.ResponseWriter, r *http.Request, cmdChan chan []byte, bCastGroup *bcast.Group) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	bCastMember :=  bCastGroup.Join()

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, bCastMember, stdoutDone)
	go ping(ws, stdoutDone)

	pumpStdin(ws, cmdChan)

	select {
	case <-stdoutDone:
	}
}

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
	http.ServeFile(w, r, "internal/webui/www/"+r.URL.Path)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		serveFile(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "internal/webui/www/index.html")
}

// StartWeb running webserver and handle pages.
func StartWeb(config_ Config, cmdChan chan []byte, outChan chan []byte) {
	bCastGroup := bcast.NewGroup() // create broadcast bcastGroup
	go bCastGroup.Broadcast(0)     // accepts messages and broadcast it to all members
	go func() {
		for {
			msg := <- outChan
			bCastGroup.Send(msg)
		}
	}()

	config = config_
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, cmdChan, bCastGroup)
	})
	log.Fatal(http.ListenAndServe(config.ListenIP+":"+config.ListenPort, nil))
}
