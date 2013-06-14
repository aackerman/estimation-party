package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/aackerman/guid"
	"log"
	"strconv"
	"time"
)

type EstimationParty struct {
	Rooms []*Room
}

func CreateRoom() Room {
	uuid, err := guid.Generate()

	if err != nil {
		log.Println("error creating guid", err)
	}

	return Room{
		Guid:    uuid,
		Voters:  make([]*Voter, 10),
		Results: Msg{Route: "results", Data: make(map[string]string)},
		Voting:  false,
		Ticket:  "",
		done:    make(chan bool),
	}
}

func SendMsg(ws *websocket.Conn, msg Msg) {
	if err := websocket.JSON.Send(ws, &msg); err != nil {
		log.Println("send err", err)
	}
}

var party = &EstimationParty{}

func (party *EstimationParty) CastVote(voter *Voter, vote Vote) {
	// handle string <-> int conversion
	i, _ := strconv.Atoi(party.Results.Data[vote.Points])
	party.Results.Data[vote.Points] = strconv.Itoa(i + 1)
	voter.Voted = true
}

func (party *EstimationParty) RemoveVoter(v *Voter) {
	i := party.FindVoter(v)
	// TODO: figure out how this works https://code.google.com/p/go-wiki/wiki/SliceTricks
	party.Voters = append(party.Voters[:i], party.Voters[i+1:]...)
}

func (party *EstimationParty) CheckVoting() {
	log.Println("CheckVoting called")
	votes := 0
	voters := 0
	for _, voter := range party.Voters {
		if voter != nil && voter.CanVote {
			voters += 1
			if voter.Voted == true {
				votes += 1
			}
		}
	}
	if votes == voters {
		log.Println("Ending estimation early")
		party.done <- true
	}
}

func (party *EstimationParty) SendResults() {
	log.Println("SendResults called to voters")
	for _, voter := range party.Voters {
		if voter != nil {
			SendMsg(voter.ws, party.Results)
			voter.Voted = false
			voter.CanVote = false
		}
	}
	party.Reset()
	log.Println("Reset Party, waiting to start estimating again")
}

func (party *EstimationParty) MakeVoter(ws *websocket.Conn) Voter {
	return Voter{
		ws:      ws,
		Voted:   false,
		CanVote: false,
		quit:    make(chan bool),
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
	log.Println("StartVoting called")
	party.Voting = true

	for _, voter := range party.Voters {
		if voter != nil {
			voter.CanVote = true
			SendMsg(voter.ws, Msg{
				Route: "start",
				Data:  map[string]string{"ticket": party.Ticket, "voting": "true"},
			})
		}
	}

	for {
		select {
		case <-party.done:
			log.Println("Estimation done early!")
			party.Voting = false
			party.SendResults()
			return
		case <-time.After(5 * time.Minute):
			log.Println("Estimation Timed Out!")
			party.Voting = false
			party.SendResults()
			return
		}
	}
}

func (party *EstimationParty) Reset() {
	party.Results = Msg{Route: "results", Data: make(map[string]string)}
	party.Ticket = ""
}
