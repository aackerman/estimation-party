package main

import (
	"code.google.com/p/go.net/websocket"
)

type Vote struct {
	Points string
}

type Voter struct {
	ws      *websocket.Conn
	Voted   bool
	CanVote bool
	quit    chan bool
}
