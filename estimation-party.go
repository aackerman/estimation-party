package main

import (
	"flag"
	"log"
	"net/http"
	// "sockpuppet"
	"encoding/json"
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
	values map[float32]int
}

type Vote struct {
	Points float32 `json:",string"`
}

var results = Results{make(map[float32]int)}

func main() {
	flag.Parse()

	// setup index handler
	http.HandleFunc("/", index)

	// listen for websocket
	SockPuppet := Listen()

	// setup websocket routing
	SockPuppet.Routing(func(s *socket) {
		var v Vote

		s.On("vote", func(data *json.RawMessage) {
			json.Unmarshal(*data, &v)
			results.values[v.Points] += 1
			log.Println(results)
		})

	})

	// server static files
	mux.Handle("/", http.FileServer(http.Dir("/Users/aackerman/www/estimation-party/public")))

	// set http server to listen
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
