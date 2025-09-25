package server

import (
	"encoding/json"
	"fmt"
	"log"
)

func HandleRoomCreate(hub *Hub, player *Player) {
	room := hub.HubCreateRoom(player)
	player.Room = room

	resp := CreateRoomResp{
		Type:        "roomCreated",
		RoomCreated: true,
		RoomCode:    room.Code,
		Host:        true,
	}
	fmt.Println("sent:", resp)
	err := player.Conn.WriteJSON(resp)
	if err != nil {
		log.Println("error sending create room message:", err)
	}
	resp2 := GameStateResp{
		Type:      "gameState",
		GameState: *room.GameState,
	}
	err2 := player.Conn.WriteJSON(resp2)
	if err2 != nil {
		log.Println("error sending Game State from create room:", err2)
	}
}

func HandleRoomJoin(hub *Hub, player *Player, msg []byte) {
	hub.Lock()
	defer hub.Unlock()
	var joinReq JoinRoomRequest
	if err := json.Unmarshal(msg, &joinReq); err != nil {
		log.Println("error unmarshaling json:", err)
	}

	fmt.Println(joinReq.RoomCode)
	room := hub.HubFindRoom(joinReq.RoomCode)
	if room == nil {
		resp := JoinRoomResp{
			Type:   "joinedRoom",
			Joined: false,
		}

		fmt.Println("sent:", resp)
		err := player.Conn.WriteJSON(resp)
		if err != nil {
			log.Println("error sending join room message:", err)
		}
		return
	}

	room.Lobby = append(room.Lobby, player)
	player.Room = room

	resp := JoinRoomResp{
		Type:   "joinedRoom",
		Joined: true,
	}

	fmt.Println("sent:", resp)
	err := player.Conn.WriteJSON(resp)
	if err != nil {
		log.Println("error sending join room message:", err)
	}
}

func HandleTeamJoin(hub *Hub, player *Player, msg []byte) {
	var joinTeamReq JoinTeamRequest
	if err := json.Unmarshal(msg, &joinTeamReq); err != nil {
		log.Println("error unmarshaling json:", err)
	}

	if joinTeamReq.Team == "left" {
		player.Room.LeftTeam = append(player.Room.LeftTeam, player)
		resp := JoinTeamResp{
			Type:   "joinedTeam",
			Joined: true,
			Team:   "left",
		}
		err := player.Conn.WriteJSON(resp)
		if err != nil {
			log.Println("error sending team joined message:", err)
		}
		return
	}
	if joinTeamReq.Team == "right" {
		player.Room.RightTeam = append(player.Room.RightTeam, player)
		resp := JoinTeamResp{
			Type:   "joinedTeam",
			Joined: true,
			Team:   "right",
		}
		err := player.Conn.WriteJSON(resp)
		if err != nil {
			log.Println("error sending team joined message:", err)
		}
		return
	}
}

func HandleRoomLeave(hub *Hub, player *Player) {
	hub.Lock()
	defer hub.Unlock()

	room := player.Room
	if room == nil {
		return
	}

	room.removePlayer(player)
	if !room.isEmpty() {
		if room.Host == player {
			room.Host = room.Lobby[0]
			resp := HostReassignment{
				Type:     "hostReassigned",
				Host:     true,
				RoomCode: room.Code,
			}
			err := room.Lobby[0].Conn.WriteJSON(resp)
			if err != nil {
				log.Println("error sending reassinged host message:", err)
			}
		}
	}
	if room.isEmpty() {
		if room.gameRunning {
			fmt.Println("rooms about to be done")
			room.done <- true
		}
		hub.HubDeleteRoom(room)
		fmt.Println("Room Deleted Because Empty")
	}

	player.Room = nil

	resp := LeaveRoomResp{
		Type:     "roomLeft",
		LeftRoom: true,
	}

	fmt.Println("sent:", resp)
	err := player.Conn.WriteJSON(resp)
	if err != nil {
		log.Println("error sending create room message:", err)
	}
}

func HandleGameStart(hub *Hub, player *Player) {
	if player.Room != nil && player.Room.Host == player && !player.Room.gameRunning {
		player.Room.gameRunning = true
		go player.Room.gameLoop()
		log.Println("go routine started for ", player.Room.Code)
	}
}

func HandleMoveUpdate(hub *Hub, player *Player, msg []byte) {
	var move MoveVoteRequest
	err := json.Unmarshal(msg, &move)
	if err != nil {
		log.Println("error with unmarshaling json:", err)
	}

	player.move = move.Direction
}

func HandleGameToLobby(hub *Hub, player *Player) {
	player.Room.done <- true
	player.Room.gameRunning = false

	player.Room.LeftTeam = player.Room.LeftTeam[:0]
	player.Room.RightTeam = player.Room.RightTeam[:0]

	resp := ReturnToLobby{
		Type: "returnToLobby",
	}

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
	gameStateResp := GameStateResp{
		Type:      "gameState",
		GameState: gameState,
	}

	for i := 0; i < len(player.Room.Lobby); i++ {
		err := player.Room.Lobby[i].Conn.WriteJSON(resp)
		if err != nil {
			log.Println("error broadcasting gamestate:", err)
		}

		err1 := player.Room.Lobby[i].Conn.WriteJSON(gameStateResp)
		if err1 != nil {
			log.Println("error broadcasting gamestate:", err1)
		}
	}
}
