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
	player.SendJSON(resp)

	resp2 := GameStateResp{
		Type:      "gameState",
		GameState: *room.GameState,
	}
	player.SendJSON(resp2)
}

func HandleRoomJoin(hub *Hub, player *Player, msg []byte) {
	hub.Lock()
	defer hub.Unlock()
	var joinReq JoinRoomRequest
	if err := json.Unmarshal(msg, &joinReq); err != nil {
		log.Println("error unmarshaling json:", err)
		return
	}

	fmt.Println(joinReq.RoomCode)
	room := hub.HubFindRoom(joinReq.RoomCode)
	if room == nil {
		resp := JoinRoomResp{
			Type:     "joinedRoom",
			Joined:   false,
			RoomCode: player.Room.Code,
		}

		fmt.Println("sent:", resp)
		player.SendJSON(resp)
		return
	}

	room.Lobby = append(room.Lobby, player)
	player.Room = room

	resp := JoinRoomResp{
		Type:     "joinedRoom",
		Joined:   true,
		RoomCode: player.Room.Code,
	}

	fmt.Println("sent:", resp)
	player.SendJSON(resp)
}

func HandleTeamJoin(hub *Hub, player *Player, msg []byte) {
	var joinTeamReq JoinTeamRequest
	if err := json.Unmarshal(msg, &joinTeamReq); err != nil {
		log.Println("error unmarshaling json:", err)
		return
	}

	if joinTeamReq.Team == "left" {
		player.Room.LeftTeam = append(player.Room.LeftTeam, player)
		resp := JoinTeamResp{
			Type:   "joinedTeam",
			Joined: true,
			Team:   "left",
		}
		player.SendJSON(resp)
		return
	}
	if joinTeamReq.Team == "right" {
		player.Room.RightTeam = append(player.Room.RightTeam, player)
		resp := JoinTeamResp{
			Type:   "joinedTeam",
			Joined: true,
			Team:   "right",
		}
		player.SendJSON(resp)
		return
	}
}

func HandleRoomLeave(hub *Hub, player *Player) {
	HandleRoomLeaveWithNotification(hub, player, true)
}

func HandleRoomLeaveWithNotification(hub *Hub, player *Player, sendNotification bool) {
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
			room.Lobby[0].SendJSON(resp)
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

	// Only send notification if player is still connected
	if sendNotification {
		resp := LeaveRoomResp{
			Type:     "roomLeft",
			LeftRoom: true,
		}

		fmt.Println("sent:", resp)
		player.SendJSON(resp)
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
	if err := json.Unmarshal(msg, &move); err != nil {
		log.Println("error with unmarshaling json:", err)
		return
	}

	player.move = move.Direction
}

func HandleGameToLobby(hub *Hub, player *Player) {
	if player.Room.gameRunning {
		player.Room.done <- true
		player.Room.gameRunning = false
	}

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

	player.Room.GameState = &gameState
	gameStateResp := GameStateResp{
		Type:      "gameState",
		GameState: gameState,
	}

	for i := 0; i < len(player.Room.Lobby); i++ {
		player.Room.Lobby[i].SendJSON(resp)
		player.Room.Lobby[i].SendJSON(gameStateResp)
	}
}
