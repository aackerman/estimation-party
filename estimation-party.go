package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"log"
	"net/http"
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

type Results struct {
	Values map[float32]int
}

type Vote struct {
	Points float32 `json:",string"`
}

type Voter struct {
	Vote Vote
}

type State struct {
	Voters  []*Voter
	Results Results
}

var state = &State{
	Voters:  make([]*Voter, 10),
	Results: Results{make(map[float32]int)},
}

func Connect(ws *websocket.Conn) {
	log.Println("socket connection established")
	defer ws.Close()

	var msg map[string]*json.RawMessage
	for {
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			log.Println("receive error", err)
			break
		}
		Route(msg)
	}

	log.Println("socket disconnected")
}

func Send(ws *websocket.Conn, msg map[string]string) {
	if err := websocket.JSON.Send(ws, &msg); err != nil {
		log.Println("send err", err)
	}
}

func Route(msg map[string]*json.RawMessage) {
	var route string
	json.Unmarshal(*msg["route"], &route)
	switch route {
	case "vote":
		var vote Vote
		json.Unmarshal(*msg["data"], &vote)
		CastVote(vote)
	case "start"
		StartVoting()
	}
}

func CastVote(v Vote) {
	log.Println(v)
}

func StartVoting() {

}

func main() {
	flag.Parse()

	// setup index handler
	http.HandleFunc("/", index)

	// listen for websocket
	http.Handle("/ws", websocket.Handler(Connect))

	// server static files
	mux.Handle("/", http.FileServer(http.Dir("/Users/aackerman/www/estimation-party/public")))

	// set http server to listen
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
