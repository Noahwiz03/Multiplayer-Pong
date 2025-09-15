package server

import (
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn *websocket.Conn
	ID   string
	Room *Room
}

type Room struct {
	Code      string
	LeftTeam  []*Player
	RightTeam []*Player
	Lobby     []*Player
	Host      *Player
	GameState *GameState
}

type GameState struct {
	LeftPaddle  Paddle
	RightPaddle Paddle
	Ball        Ball
	ScoreLeft   int
	ScoreRight  int
}

type Paddle struct {
	X, Y   float64
	Width  float64
	Height float64
	Speed  float64
}

type Ball struct {
	X, Y      float64
	Radius    float64
	VelocityX float64
	VelocityY float64
}

func CreateRoom(player *Player) *Room {
	roomCode := generateRoomCode()

	gameState := GameState{
		LeftPaddle: Paddle{
			X:      30,
			Y:      250,
			Width:  20,
			Height: 100,
			Speed:  5,
		},
		RightPaddle: Paddle{
			X:      750,
			Y:      250,
			Width:  20,
			Height: 100,
			Speed:  5,
		},
		Ball: Ball{
			X:         400,
			Y:         300,
			Radius:    10,
			VelocityX: 0,
			VelocityY: 0,
		},
		ScoreLeft:  0,
		ScoreRight: 0,
	}

	room := &Room{
		Code:      roomCode,
		LeftTeam:  []*Player{},
		RightTeam: []*Player{},
		Lobby:     []*Player{player},
		Host:      player,
		GameState: &gameState,
	}
	return room
}

func (r *Room) removePlayer(player *Player) {
	r.LeftTeam = removePlayerFromSlice(r.LeftTeam, player)
	r.RightTeam = removePlayerFromSlice(r.RightTeam, player)
	r.Lobby = removePlayerFromSlice(r.Lobby, player)
}

func removePlayerFromSlice(slice []*Player, player *Player) []*Player {
	for i, p := range slice {
		if p == player {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (r *Room) isEmpty() bool {
	return (len(r.LeftTeam) == 0 && len(r.RightTeam) == 0 && len(r.Lobby) == 0)
}

func generateRoomCode() string {
	return strings.ToUpper(uuid.New().String()[:6]) // e.g., "A1B2C3"
}

//gamestate related things
