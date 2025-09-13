package server

import (
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn *websocket.Conn
	ID   string
	RoomCode string
}

type Room struct {
	Code      string
	LeftTeam  []*Player
	RightTeam []*Player
	Lobby     []*Player
}

func CreateRoom(player *Player) *Room{
	roomCode := generateRoomCode()

	room := &Room{
		Code: roomCode,
		LeftTeam: []*Player{},
		RightTeam: []*Player{},
		Lobby: []*Player{player},
	}
	return room
}

func (r *Room) removePlayer(player *Player){
	r.LeftTeam = removePlayerFromSlice(r.LeftTeam, player)
	r.RightTeam = removePlayerFromSlice(r.RightTeam, player)
	r.Lobby = removePlayerFromSlice(r.Lobby, player)
}

func removePlayerFromSlice(slice []*Player, player *Player) []*Player{
	for i, p := range slice{
		if p == player{
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (r *Room)isEmpty() bool{
	return (len(r.LeftTeam) == 0 && len(r.RightTeam) == 0 && len(r.Lobby) == 0)
}

func generateRoomCode() string {
	return strings.ToUpper(uuid.New().String()[:6]) // e.g., "A1B2C3"
}
