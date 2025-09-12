package server

import (
	"sync"
)

type Hub struct {
	Rooms map[string]*Room
	sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub)HubCreateRoom(player *Player) *Room{
	room := CreateRoom(player)
	h.Lock()
	h.Rooms[room.Code] = room
	h.Unlock()
	return room
	//add room to rooms map
}
