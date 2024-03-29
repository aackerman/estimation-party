package main

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"github.com/aackerman/guid"
	"io"
	"log"
	"strconv"
	"time"
)

type Room struct {
	Guid    string
	Voters  map[*Voter]bool
	Voting  bool
	Ticket  string
	Results Msg
	done    chan bool
}

func CreateRoom() *Room {
	uuid, err := guid.Generate()

	if err != nil {
		log.Println("error creating guid", err)
	}

	log.Println("Created a new room at", uuid)
	room := &Room{
		Guid:    uuid,
		Voters:  make(map[*Voter]bool, 10),
		Results: Msg{Route: "results", Data: make(map[string]string)},
		Voting:  false,
		Ticket:  "",
		done:    make(chan bool),
	}

	EstimationParty.Rooms[room] = true

	return room
}

func (r *Room) Listen(v *Voter) {
	var msg Msg

	r.Voters[v] = true

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

func (r *Room) CastVote(voter *Voter, vote Vote) {
	// handle string <-> int conversion
	i, _ := strconv.Atoi(r.Results.Data[vote.Points])
	r.Results.Data[vote.Points] = strconv.Itoa(i + 1)
	voter.Voted = true
}

func (r *Room) RemoveVoter(v *Voter) {
	delete(r.Voters, v)
}

func (r *Room) CheckVoting() {
	log.Println("Checking voting to end early")
	votes := 0
	voters := 0
	for voter, _ := range r.Voters {
		if voter.CanVote == true {
			voters += 1
			if voter.Voted == true {
				votes += 1
			}
		}
	}
	if votes == voters {
		log.Println("Ending estimation early")
		r.done <- true
	}
}

func (r *Room) SendResults() {
	log.Println("Sending results to voters")
	for voter, _ := range r.Voters {
		r.Results.Send(voter.ws)
		voter.Voted = false
		voter.CanVote = false
	}
	r.Reset()
}

func (r *Room) StartVoting() {
	r.Voting = true

	log.Println("Sent socket messages to voters to start voting")
	for voter, _ := range r.Voters {
		voter.CanVote = true
		msg := &Msg{
			Route: "start",
			Data:  map[string]string{"ticket": r.Ticket, "voting": "true"},
		}
		msg.Send(voter.ws)
	}

	for {
		select {
		case <-r.done:
			log.Println("Estimation done early!")
			r.Voting = false
			r.SendResults()
			return
		case <-time.After(5 * time.Minute):
			log.Println("Estimation Timed Out!")
			r.Voting = false
			r.SendResults()
			return
		}
	}
}

func FindRoomByGuid(guid string) (*Room, error) {
	log.Println("Room we want:", guid)
	for room, _ := range EstimationParty.Rooms {
		log.Println("Room to check:", room.Guid)
		log.Println("match?", guid == room.Guid)
		if room.Guid == guid {
			return room, nil
		}
	}
	return &Room{}, errors.New("Room does not exist")
}

func (r *Room) Reset() {
	log.Println("Reset Party, waiting to start estimating again")
	r.Results = Msg{Route: "results", Data: make(map[string]string)}
	r.Ticket = ""
}
