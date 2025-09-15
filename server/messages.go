package server

// requests
type CreateRoomRequest struct {
	Type string `json:"type"`
}
type JoinRoomRequest struct {
	Type     string `json:"type"`
	RoomCode string `json:"roomCode"`
}
type JoinTeamRequest struct {
	Type string `json:"type"`
	Team string `json:"team"`
}
type LeaveRoomRequest struct {
	Type string `json:"type"`
}

// responses
type CreateRoomResp struct {
	Type        string `json:"type"`
	RoomCreated bool   `json:"roomCreated"`
	RoomCode    string `json:"roomCode"`
	Host        bool   `json:"host"`
}
type JoinRoomResp struct {
	Type   string `json:"type"`
	Joined bool   `json:"joined"`
}
type JoinTeamResp struct {
	Type   string `json:"type"`
	Joined bool   `json:"joined"`
	Team   string `json:"team"`
}
type LeaveRoomResp struct {
	Type     string `json:"type"`
	LeftRoom bool   `json:"leftRoom"`
	Host     bool   `json:"host"`
}
type HostReassignment struct {
	Type string `json:"type"`
	Host bool   `json:"host"`
}
