package server

import(
	"log"
	"fmt"
)

func HandleRoomCreate(hub *Hub, player *Player){
	room := hub.HubCreateRoom(player)

	player.RoomCode = room.Code

	resp := CreateRoomResp{
		Type: "roomCreated",
		RoomCreated: true,
		RoomCode: room.Code,
	}

	fmt.Println("sent:",resp)
	err:= player.Conn.WriteJSON(resp)
	if err != nil{
		log.Println("error sending create room message:", err)
	}
}

func HandleRoomLeave(hub *Hub, player *Player){
	hub.Lock()
	defer hub.Unlock()

	room, exists := hub.Rooms[player.RoomCode]
	if !exists{
		return
	}

	room.removePlayer(player)

	if room.isEmpty() {
		hub.HubDeleteRoom(room)
		fmt.Println("Room Deleted Because Empty")
	}

	player.RoomCode = ""

	resp := LeaveRoomResp{
		Type: "roomLeft",
		LeftRoom: true,	
	} 

	fmt.Println("sent:",resp)
	err:= player.Conn.WriteJSON(resp)
	if err != nil{
		log.Println("error sending create room message:", err)
	}
}

