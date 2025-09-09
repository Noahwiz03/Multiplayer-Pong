package server

import (
	"sync"
)

type Hub struct {
	Rooms map[string]*Room
	sync.Mutex
}
