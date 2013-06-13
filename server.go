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

func main() {
	flag.Parse()

	// setup index handler
	http.HandleFunc("/", index)

	// listen for websocket
	http.Handle("/ws", websocket.Handler(Connect))

	dirs := []string{os.Getenv("HOME"), "/www/estimation-party/public"}
	dir := strings.Join(dirs, "")

	// server static files
	mux.Handle("/", http.FileServer(http.Dir(dir)))

	// set http server to listen
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
