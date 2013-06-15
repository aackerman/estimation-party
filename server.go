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
		indextpl.Execute(w, "")
	}
}

func roomHandler(w http.ResponseWriter, r *http.Request) {
	subpath := r.URL.Path[6:]

	if subpath == "new" {
		room := CreateRoom()
		url := "/room/" + room.Guid
		http.Redirect(w, r, url, 302)
		return
	}

	room, err := FindRoomByGuid(subpath)

	if err != nil {
		http.Redirect(w, r, "/", 404)
		return
	}

	roomtpl.Execute(w, room)
}

func WebsocketConnect(ws *websocket.Conn) {
	defer ws.Close()

	voter := MakeVoter(ws)

	guid := ws.Request().URL.Path[4:]
	room, err := FindRoomByGuid(guid)

	if err != nil {
		log.Fatal("Room does not exist")
	}

	go room.Listen(voter)

	<-voter.quit
	delete(room.Voters, voter)
}

func main() {
	flag.Parse()

	// setup index handler
	http.HandleFunc("/", index)

	http.HandleFunc("/room/", roomHandler)

	// listen for websockets
	http.Handle("/ws/", websocket.Handler(WebsocketConnect))

	dirs := []string{os.Getenv("HOME"), "/www/estimation-party/public"}
	dir := strings.Join(dirs, "")

	// server static files
	mux.Handle("/", http.FileServer(http.Dir(dir)))

	// set http server to listen
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
