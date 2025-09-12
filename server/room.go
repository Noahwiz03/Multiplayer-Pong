package server

import (
	"github.com/gorilla/websocket"
	"strings"
	"github.com/google/uuid"
)

type Player struct {
	Conn *websocket.Conn
	ID   string
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

func generateRoomCode() string {
	return strings.ToUpper(uuid.New().String()[:6]) // e.g., "A1B2C3"
}
