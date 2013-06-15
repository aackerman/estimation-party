package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var tpl = template.Must(template.ParseFiles("public/index.html"))
var mux = http.NewServeMux()

func index(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		mux.ServeHTTP(w, r)
	} else {
		tpl.Execute(w, r.Host)
	}
}

func roomHandler(w http.ResponseWriter, r *http.Request) {
	subpath := r.URL.Path[6:]
	switch subpath {
	case "new":
		// create room and redirect user to the room
	case ExistingRoom(subpath):
		// respond with the template for the room
	default:
		// room does not exist, redirect to 404
	}
}

func ExistingRoom(name string) string {
	return name
	// check if the room is in map of existing rooms
}

func FindRoom(ws *websocket.Conn) Room {
	return Room{}
}

func WebsocketConnect(ws *websocket.Conn) {
	log.Println("socket connected")
	defer ws.Close()

	room := FindRoom(ws)

	// make a new voter
	voter := room.MakeVoter(ws)
	room.Voters[voter] = true

	go room.Receive(voter)
	<-voter.quit
	delete(room.Voters, voter)
	log.Println("socket disconnected")
}

func main() {
	flag.Parse()

	// setup index handler
	http.HandleFunc("/", index)

	http.HandleFunc("/room/", roomHandler)

	// listen for websockets
	http.Handle("/ws", websocket.Handler(WebsocketConnect))

	dirs := []string{os.Getenv("HOME"), "/www/estimation-party/public"}
	dir := strings.Join(dirs, "")

	// server static files
	mux.Handle("/", http.FileServer(http.Dir(dir)))

	// set http server to listen
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
