package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"strconv"
	"time"
)

type EstimationParty struct {
	Voters  []*Voter
	Voting  bool
	Ticket  string
	Results Msg
	done    chan bool
}

type Msg struct {
	Route string            `json:"route"`
	Data  map[string]string `json:"data"`
}

var party = &EstimationParty{
	Voters:  make([]*Voter, 10),
	Results: Msg{Route: "results", Data: make(map[string]string)},
	Voting:  false,
	Ticket:  "",
	done:    make(chan bool),
}

func (party *EstimationParty) CastVote(voter *Voter, vote Vote) {
	i, _ := strconv.Atoi(party.Results.Data[vote.Points])
	party.Results.Data[vote.Points] = strconv.Itoa(i + 1)
	voter.Voted = true
}

func (party *EstimationParty) RemoveVoter(v *Voter) {
	i := party.FindVoter(v)
	party.Voters = append(party.Voters[:i], party.Voters[i+1:]...)
}

func (party *EstimationParty) CheckVoting() {
	log.Println("CheckVoting called")
	votes := 0
	voters := 0
	for _, voter := range party.Voters {
		if voter != nil {
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
			voter.SendMsg(party.Results)
			voter.Voted = false
		}
	}
	party.Reset()
}

func (party *EstimationParty) MakeVoter(ws *websocket.Conn) Voter {
	return Voter{
		ws:    ws,
		Voted: false,
		quit:  make(chan bool),
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
			voter.SendMsg(Msg{
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
		case <-time.After(20 * time.Second):
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

func Connect(ws *websocket.Conn) {
	log.Println("socket connected")
	defer ws.Close()

	// make a new voter
	voter := party.MakeVoter(ws)

	// append voter to the app state
	party.Voters = append(party.Voters, &voter)

	go voter.Receive()

	<-voter.quit

	// remove voter from state.Voters
	party.RemoveVoter(&voter)
	log.Println("socket disconnected")
}
