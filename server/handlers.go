package server

import(
	"log"
	"fmt"
)

func HandleRoomCreate(hub *Hub, player *Player){
	room := hub.HubCreateRoom(player)

	resp := CreateRoomResp{
		Type: "roomCreated",
		RoomCreated: true,
		RoomCode: room.Code,
	}

	fmt.Println(resp)
	err:= player.Conn.WriteJSON(resp)
	if err != nil{
		log.Println("error sending create room message:", err)
	}
}


