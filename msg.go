package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type Msg struct {
	Route string            `json:"route"`
	Data  map[string]string `json:"data"`
}

func (m *Msg) Send(ws *websocket.Conn) {
	if err := websocket.JSON.Send(ws, &m); err != nil {
		log.Println("send err", err)
	}
}
