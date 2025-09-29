package server

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn *websocket.Conn
	ID   string
	Room *Room
	move int //1 up 0 no move -1 down
	send chan []byte
	done chan struct{}
}

type Room struct {
	Code        string
	LeftTeam    []*Player
	RightTeam   []*Player
	Lobby       []*Player
	Host        *Player
	GameState   *GameState
	done        chan bool
	mutex       sync.RWMutex
	gameRunning bool
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
			X:      30,  //top left
			Y:      250, //top left
			Width:  10,
			Height: 100,
			Speed:  5,
		},
		RightPaddle: Paddle{
			X:      750,
			Y:      250,
			Width:  10,
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
		Code:        roomCode,
		LeftTeam:    []*Player{},
		RightTeam:   []*Player{},
		Lobby:       []*Player{player},
		Host:        player,
		GameState:   &gameState,
		done:        make(chan bool, 1),
		gameRunning: false,
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

// gamestate related things
// such as getting move requests! and sending out gamestate updates
func (r *Room) gameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()
	r.GameState.Ball.VelocityX = 4
	r.GameState.Ball.VelocityY = 2

	for {
		select {
		case <-ticker.C:
			r.updateGameState()
			r.broadcastGameState()
		case <-r.done:
			return
		}
	}
}

func (r *Room) updateGameState() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	leftMoveTally := 0
	rightMoveTally := 0

	//ball collision with paddles
	if ballPaddleCollision(int(r.GameState.Ball.X), int(r.GameState.Ball.Y), int(r.GameState.Ball.Radius),
		int(r.GameState.LeftPaddle.X), int(r.GameState.LeftPaddle.Y), int(r.GameState.LeftPaddle.Width),
		int(r.GameState.LeftPaddle.Height)) {
		r.GameState.Ball.VelocityX *= -1
	}
	if ballPaddleCollision(int(r.GameState.Ball.X), int(r.GameState.Ball.Y), int(r.GameState.Ball.Radius),
		int(r.GameState.RightPaddle.X), int(r.GameState.RightPaddle.Y), int(r.GameState.RightPaddle.Width),
		int(r.GameState.RightPaddle.Height)) {
		r.GameState.Ball.VelocityX *= -1
	}

	//ball collision with top and bottom
	if r.GameState.Ball.Y-r.GameState.Ball.Radius <= 0 {
		r.GameState.Ball.VelocityY *= -1
	}
	if r.GameState.Ball.Y+r.GameState.Ball.Radius >= 600 {
		r.GameState.Ball.VelocityY *= -1
	}

	//ball scoreing
	if r.GameState.Ball.X-r.GameState.Ball.Radius <= 0 {
		r.GameState.ScoreRight++
		r.GameState.Ball.X = 400
		r.GameState.Ball.Y = 300
	}
	if r.GameState.Ball.X+r.GameState.Ball.Radius >= 800 {
		r.GameState.ScoreLeft++
		r.GameState.Ball.X = 400
		r.GameState.Ball.Y = 300
	}

	r.GameState.Ball.X += r.GameState.Ball.VelocityX
	r.GameState.Ball.Y += r.GameState.Ball.VelocityY

	for i := 0; i < len(r.RightTeam); i++ {
		rightMoveTally += r.RightTeam[i].move
	}

	for i := 0; i < len(r.LeftTeam); i++ {
		leftMoveTally += r.LeftTeam[i].move
	}

	//make the paddle move
	if rightMoveTally < 0 && r.GameState.RightPaddle.Y+100 < 600 {
		r.GameState.RightPaddle.Y += r.GameState.RightPaddle.Speed
	}
	if rightMoveTally > 0 && r.GameState.RightPaddle.Y > 0 {
		r.GameState.RightPaddle.Y -= r.GameState.RightPaddle.Speed
	}

	if leftMoveTally < 0 && r.GameState.LeftPaddle.Y+100 < 600 {
		r.GameState.LeftPaddle.Y += r.GameState.LeftPaddle.Speed
	}
	if leftMoveTally > 0 && r.GameState.LeftPaddle.Y > 0 {
		r.GameState.LeftPaddle.Y -= r.GameState.LeftPaddle.Speed
	}

}

func ballPaddleCollision(bx int, by int, r int, rx int, ry int, rw int, rh int) bool {
	closestX := max(rx, min(bx, rx+rw))
	closestY := max(ry, min(by, ry+rh))

	dx := bx - closestX
	dy := by - closestY

	return (dx*dx + dy*dy) <= r*r
}

func (r *Room) broadcastGameState() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	resp := GameStateResp{
		Type:      "gameState",
		GameState: *r.GameState,
	}
	for i := 0; i < len(r.Lobby); i++ {
		err := r.Lobby[i].SendJSON(resp)
		if err != nil {
			log.Println("error broadcasting gamestate:", err)
		}
	}
}
