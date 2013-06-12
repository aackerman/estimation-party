package main

import (
	"flag"
	"log"
	"net/http"
	// "sockpuppet"
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

	IO := Listen()

	IO.Sockets(func(s *socket) {

		s.On("hello", func(msg map[string]string) {
			s.Send(msg)
			log.Println("echo", msg)
		})

	})

	// server static files
	mux.Handle("/", http.FileServer(http.Dir("/Users/aackerman/www/estimation-party/public")))

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
