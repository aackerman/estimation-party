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
	ws      *websocket.Conn
	Voted   bool
	CanVote bool
	quit    chan bool
}

func (r *Room) Receive(v *Voter) {
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
			if r.Voting == true && v.Voted == false && v.CanVote {
				r.CastVote(v, vote)
				r.CheckVoting()
			}
		case "start":
			r.Ticket = msg.Data["ticket"]
			go r.StartVoting()
		default:
			log.Println("unknown route", msg.Route)
		}
	}
}
