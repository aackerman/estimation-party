package main

import (
	"code.google.com/p/go.net/websocket"
	_ "encoding/json"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type EstimationParty struct {
	Voters  []*Voter
	Voting  bool
	Ticket  string
	Results map[float64]int
	done    chan bool
}

type Vote struct {
	Points float64 `json:",string"`
}

type Voter struct {
	ws       *websocket.Conn
	id       int
	Vote     Vote
	receive  chan Msg
	response chan Msg
	quit     chan bool
}

type Msg struct {
	Route string
	Data  map[string]string
}

var party = &EstimationParty{
	Voters:  make([]*Voter, 10),
	Results: make(map[float64]int),
}

func (v *Voter) Send() {
	for {
		select {
		case msg := <-v.response:
			if err := websocket.JSON.Send(v.ws, &msg); err != nil {
				log.Println("send err", err)
				break
			}
		}
	}
}

func (v *Voter) Receive() {
	var msg Msg
	var route string
	var vote Vote
	var points float64

	for {
		if err := websocket.JSON.Receive(v.ws, &msg); err != nil {
			if err != io.EOF {
				log.Println("websocket receive error", err)
			}
			v.quit <- true
			return
		}

		switch route {
		case "vote":
			points = strconv.ParseFloat(msg.Data["points"], 64)
			vote = &Vote{points}
			if party.Voting == true {
				party.CastVote(v, vote)
			}
		case "start":
			log.Println("start voting!")
			go party.StartVoting()
		// case "sync":
		default:
			v.quit <- true
			return
		}
	}
}

func (party *EstimationParty) CastVote(voter *Voter, vote Vote) {
	party.Results[vote.Points] += 1
	voter.Vote = vote
}

func (party *EstimationParty) RemoveVoter(v *Voter) {
	i := party.FindVoter(v)
	party.Voters = append(party.Voters[:i], party.Voters[i+1:]...)
}

func (party *EstimationParty) MakeVoter(ws *websocket.Conn) Voter {
	return Voter{
		ws:       ws,
		id:       rand.Int(),
		response: make(chan Msg),
		receive:  make(chan Msg),
		quit:     make(chan bool),
	}
}

func (party *EstimationParty) FindVoter(v *Voter) int {
	for i, val := range party.Voters {
		if v == val {
			return i
		}
	}
	return -1
}

func (party *EstimationParty) StartVoting() {
	party.Voting = true

	for _, voter := range party.Voters {
		if voter != nil {
			voter.response <- Msg{
				Route: "start",
				Data:  map[string]string{"ticket": party.Ticket, "voting": "true"},
			}
		}
	}

	for {
		select {
		case <-party.done:
			party.Voting = false
		case <-time.After(20 * time.Second):
			party.Voting = false
		}
	}
}

func (party *EstimationParty) Reset() {
	party.Results = make(map[float64]int)
	party.Ticket = ""
}

func Connect(ws *websocket.Conn) {
	log.Println("socket connection established")
	defer ws.Close()

	// make a new voter
	voter := party.MakeVoter(ws)

	// append voter to the app state
	party.Voters = append(party.Voters, &voter)

	go voter.Receive()
	go voter.Send()

	<-voter.quit

	// remove voter from state.Voters
	party.RemoveVoter(&voter)
	log.Println("socket disconnected")
}
