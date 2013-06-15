package main

type Party struct {
	Rooms map[*Room]bool
}

var EstimationParty = &Party{
	Rooms: make(map[*Room]bool),
}
