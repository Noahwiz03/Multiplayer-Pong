package server

import (
	"fmt"
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

func (h *Hub) HubCreateRoom(player *Player) *Room {
	room := CreateRoom(player)

	h.Lock()
	h.Rooms[room.Code] = room
	h.Unlock()

	return room
}

func (h *Hub) HubFindRoom(roomcode string) *Room {
	room, exists := h.Rooms[roomcode]
	if !exists {
		return nil
	}
	return room
}

func (h *Hub) HubDeleteRoom(room *Room) {
	fmt.Println("made it into hubdelete room")
	select {
	case room.done <- true:
		fmt.Println("sent done signal to goroutine")
	default:
		fmt.Println("goroutine already stopped or not running")
	}
	delete(h.Rooms, room.Code)
	fmt.Println("room deleted from hub")
}
