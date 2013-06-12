package main

import (
	"code.google.com/p/go.net/websocket"
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

func (s *SockPuppet) Sockets(fn func(*socket)) {
	go func() {
		sock := <-DefaultSockPuppet.connect
		fn(sock)
	}()
}

func connect(ws *websocket.Conn) {
	log.Println("socket connected")
	s := &socket{
		ws,
		make(map[string]func(map[string]string)),
		make(chan bool),
	}
	go s.Main()
	DefaultSockPuppet.connect <- s
	<-s.quit
	s.disconnect()
}

type socket struct {
	ws     *websocket.Conn
	routes map[string]func(map[string]string)
	quit   chan bool
}

func (s *socket) Main() {
	var msg map[string]string
	for {
		if err := websocket.JSON.Receive(s.ws, &msg); err != nil {
			s.quit <- true
			break
		}
		log.Println("heelo")
		log.Println(msg)
		if len(msg["route"]) > 0 && s.routes[msg["route"]] != nil {
			s.routes[msg["route"]](msg)
		} else {
			log.Println("Missing route error", msg["route"])
		}
	}
}

func (s *socket) disconnect() {
	s.ws.Close()
	log.Println("websocket disconnected")
}

func (s *socket) On(str string, fn func(map[string]string)) {
	s.routes[str] = fn
}

func (s *socket) Send(data map[string]string) {
	if err := websocket.JSON.Send(s.ws, &data); err != nil {
		s.quit <- true
	}
}
