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

var indextpl = template.Must(template.ParseFiles("views/layout.html", "views/index.html"))
var roomtpl = template.Must(template.ParseFiles("views/layout.html", "views/room.html"))
var mux = http.NewServeMux()

func index(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		mux.ServeHTTP(w, r)
	} else {
		indextpl.Execute(w, r.URL.Host)
	}
}

func roomHandler(w http.ResponseWriter, r *http.Request) {
	subpath := r.URL.Path[6:]
	switch subpath {
	case "new":
		room := CreateRoom()
		url := "/room/" + room.Guid
		http.Redirect(w, r, url, 302)
	case ExistingRoom(subpath):
		roomtpl.Execute(w, "")
	default:
		// room does not exist, redirect to 404
	}
}

func ExistingRoom(name string) string {
	return name
	// check if the room is in map of existing rooms
}

func FindRoom(guid string) *Room {
	for room, _ := range EstimationParty.Rooms {
		if room.Guid == guid {
			return room
		}
	}
	return &Room{}
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
