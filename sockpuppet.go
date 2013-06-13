package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
	"net/http"
)

type SockPuppet struct {
	sockets  map[*socket]bool
	connect  chan *socket
	callback func(*socket)
}

var DefaultSockPuppet = NewSockPuppet()

func NewSockPuppet() SockPuppet {
	return SockPuppet{
		sockets: make(map[*socket]bool),
		connect: make(chan *socket),
	}
}

func Listen() SockPuppet {
	http.Handle("/ws", websocket.Handler(connect))
	return DefaultSockPuppet
}

func (s *SockPuppet) Routing(fn func(*socket)) {
	go func() {
		for {
			select {
			case sock := <-DefaultSockPuppet.connect:
				fn(sock)
			}
		}
	}()
}

func connect(ws *websocket.Conn) {
	log.Println("socket connected")
	s := &socket{
		ws,
		make(map[string]func(*json.RawMessage)),
		make(chan bool),
	}
	go s.Main()
	DefaultSockPuppet.connect <- s
	<-s.quit
	s.disconnect()
}

type socket struct {
	ws     *websocket.Conn
	routes map[string]func(*json.RawMessage)
	quit   chan bool
}

func (s *socket) Main() {
	var route string
	var msg map[string]*json.RawMessage
	for {
		if err := websocket.JSON.Receive(s.ws, &msg); err != nil {
			log.Println("receive error", err)
			s.quit <- true
			break
		}

		if err := json.Unmarshal(*msg["route"], &route); err != nil {
			log.Fatal("JSON parse error for route")
		}

		if msg["route"] != nil && s.routes[route] != nil {
			s.routes[route](msg["data"])
		} else {
			log.Println("Missing route error", msg["route"])
		}
	}
}

func (s *socket) disconnect() {
	s.ws.Close()
	log.Println("websocket disconnected")
}

func (s *socket) On(str string, fn func(*json.RawMessage)) {
	s.routes[str] = fn
}

func (s *socket) Send(data map[string]string) {
	if err := websocket.JSON.Send(s.ws, &data); err != nil {
		log.Println("send err", err)
		s.quit <- true
	}
}
