package server

import (
	"encoding/json"
	"fmt"
	"log"
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

func HandleRoomJoin(hub *Hub, player *Player, msg []byte){
	hub.Lock()
	defer hub.Unlock()
	var joinReq JoinRoomRequest
	if err := json.Unmarshal(msg,&joinReq); err != nil{
		log.Println("error unmarshaling json:", err)
	}

  fmt.Println(joinReq.RoomCode)
	room := hub.HubFindRoom(joinReq.RoomCode)
	if room == nil {
		resp := JoinRoomResp{
			Type: "joinedRoom",
			Joined: false,	
		} 

		fmt.Println("sent:", resp)
		err := player.Conn.WriteJSON(resp)
		if err != nil{
			log.Println("error sending join room message:", err)
		}
		return
	}

	room.Lobby = append(room.Lobby, player)
	player.RoomCode = room.Code

	resp := JoinRoomResp{
		Type: "joinedRoom",
		Joined: true,
	}

	fmt.Println("sent:", resp)
	err := player.Conn.WriteJSON(resp)
	if err != nil{
		log.Println("error sending join room message:", err)
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

