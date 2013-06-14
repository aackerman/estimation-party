package main

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
)

type Vote struct {
	Points string
}

type Voter struct {
	ws    *websocket.Conn
	Voted bool
	quit  chan bool
}

func (v *Voter) SendMsg(msg Msg) {
	if err := websocket.JSON.Send(v.ws, &msg); err != nil {
		log.Println("send err", err)
	}
}

func (v *Voter) Receive() {
	var msg Msg

	for {
		if err := websocket.JSON.Receive(v.ws, &msg); err != nil {
			if err != io.EOF {
				log.Println("websocket receive error", err)
			}
			v.quit <- true
			return
		}

		switch msg.Route {
		case "vote":
			var vote Vote
			vote = Vote{msg.Data["points"]}
			if party.Voting == true && v.Voted == false {
				party.CastVote(v, vote)
				party.CheckVoting()
			}
		case "start":
			party.Ticket = msg.Data["ticket"]
			go party.StartVoting()
		default:
			log.Println("unknown route", msg.Route)
		}
	}
}
