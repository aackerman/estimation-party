package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"io"
	"log"
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
	ws           *websocket.Conn
	Voted        bool
	receive      chan Msg
	msgresponse  chan Msg
	byteresponse chan []byte
	quit         chan bool
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
		case msg := <-v.msgresponse:
			if err := websocket.JSON.Send(v.ws, &msg); err != nil {
				log.Println("send err", err)
				break
			}
		case data := <-v.byteresponse:
			if err := websocket.JSON.Send(v.ws, &data); err != nil {
				log.Println("send err", err)
				break
			}
		}
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
			points, _ := strconv.ParseFloat(msg.Data["points"], 64)
			vote = Vote{points}
			if party.Voting == true && v.Voted == false {
				party.CastVote(v, vote)
			}
		case "start":
			log.Println("start voting!")
			go party.StartVoting()
		// case "sync":
		default:
			log.Println("unknown route", msg.Route)
			v.quit <- true
			return
		}
	}
}

func (party *EstimationParty) CastVote(voter *Voter, vote Vote) {
	party.Results[vote.Points] += 1
	voter.Voted = true
}

func (party *EstimationParty) RemoveVoter(v *Voter) {
	i := party.FindVoter(v)
	party.Voters = append(party.Voters[:i], party.Voters[i+1:]...)
}

func (party *EstimationParty) CheckVoting() {
	votes := 0
	for _, voter := range party.Voters {
		if voter.Voted == true {
			votes += 1
		}
	}
	if votes >= len(party.Voters) {
		party.done <- true
	}
}

func (party *EstimationParty) ResetVotes() {
	for _, voter := range party.Voters {
		voter.Voted = false
	}
}

func (party *EstimationParty) SendResults() {
	for _, voter := range party.Voters {
		if voter != nil {
			marshaled, err := json.Marshal(party.Results)
			if err != nil {
				log.Println("json error")
			}
			voter.byteresponse <- marshaled
		}
	}
}

func (party *EstimationParty) MakeVoter(ws *websocket.Conn) Voter {
	return Voter{
		ws:           ws,
		Voted:        false,
		msgresponse:  make(chan Msg),
		byteresponse: make(chan []byte),
		receive:      make(chan Msg),
		quit:         make(chan bool),
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
			voter.msgresponse <- Msg{
				Route: "start",
				Data:  map[string]string{"ticket": party.Ticket, "voting": "true"},
			}
		}
	}

	for {
		select {
		case <-party.done:
			party.Voting = false
			party.SendResults()
			break
		case <-time.After(20 * time.Second):
			party.Voting = false
			party.SendResults()
			break
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
