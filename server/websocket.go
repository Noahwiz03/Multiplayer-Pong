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
	defer ws.Close()

	player := &Player{
		Conn: ws,
		ID:   generatePlayerID(),
		RoomCode: "",
	}

	fmt.Println("new player connected:", player.ID)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		fmt.Println("Recieved:", string(msg))
		handleMessage(hub, player, msg)
	}
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
		var req CreateRoomRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("error parsing createroomrequest:", err)
			return
		}

		HandleRoomCreate(hub, player)
		fmt.Println("room created")
	case "joinRoom":
		var req CreateRoomRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("error parsing createroomrequest:", err)
			return
		}
		fmt.Println("Room Joined")

	case "joinTeam":
		var req CreateRoomRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("error parsing createroomrequest:", err)
			return
		}
		fmt.Println("Team Joined")

	case "leaveRoom":
		var req CreateRoomRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("error parsing createroomrequest:", err)
			return
		}

		HandleRoomLeave(hub, player)
		fmt.Println("Exited Room")

	}
}

// generates a uuid for each player for unique id!
func generatePlayerID() string {
	return uuid.New().String()
}
