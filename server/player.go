package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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

// player read write pumps and JSON sending wrapper
// used AI for this but it all makes sense and it used code i already wrote for some of it
func (p *Player) readPump(hub *Hub) {
	defer func() {
		// Clean up player from room BEFORE closing connection (no notification)
		HandleRoomLeaveWithNotification(hub, p, false)
		p.Conn.Close()
		close(p.done)
	}()

	p.Conn.SetReadDeadline(time.Now().Add(70 * time.Second))

	for {
		_, msg, err := p.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		fmt.Println("Recieved:", string(msg))
		handleMessage(hub, p, msg)
	}
}

func (p *Player) writePump() {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		// Don't close connection here - readPump handles it
	}()

	for {
		select {
		case message := <-p.send:
			p.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := p.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("write error: %v", err)
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			p.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-p.done:
			// Cleanup signal received
			return
		}
	}
}

// Helper method for sending JSON messages
func (p *Player) SendJSON(v interface{}) error {
	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}

	select {
	case p.send <- msg:
		return nil
	case <-p.done:
		return errors.New("player disconnected")
	default:
		// Channel is full, player might be slow
		return errors.New("send channel full")
	}
}
