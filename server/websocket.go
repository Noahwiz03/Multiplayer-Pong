package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handles web sockets ovb
func HandleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Don't close connection here - readPump handles it

	player := &Player{
		Conn: ws,
		ID:   generatePlayerID(),
		Room: nil,
		move: 0,
		send: make(chan []byte, 256),
		done: make(chan struct{}),
	}

	fmt.Println("new player connected:", player.ID)
	go player.readPump(hub)
	go player.writePump()

}

// handles the json message based on its type
type MessageType struct {
	Type string `json:"type"`
}

func handleMessage(hub *Hub, player *Player, msg []byte) {
	var msgType MessageType
	err := json.Unmarshal(msg, &msgType)
	if err != nil {
		log.Println("Invalid JSON", err)
		return
	}

	switch msgType.Type {
	case "createRoom":
		HandleRoomCreate(hub, player)

	case "joinRoom":
		HandleRoomJoin(hub, player, msg)

	case "joinTeam":
		var req CreateRoomRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("error parsing createroomrequest:", err)
			return
		}
		HandleTeamJoin(hub, player, msg)

	case "leaveRoom":
		HandleRoomLeave(hub, player)

	case "gameStart":
		HandleGameStart(hub, player)

	case "moveUpdate":
		HandleMoveUpdate(hub, player, msg)
	case "gameToLobby":
		HandleGameToLobby(hub, player)
	}
}

// generates a uuid for each player for unique id!
func generatePlayerID() string {
	return uuid.New().String()
}
