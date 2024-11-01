package state

import (
	"sync"
)

type Clients map[string]chan string

type ClientState struct {
	Mu      sync.RWMutex
	Clients Clients
}

func NewClientState() *ClientState {
	return &ClientState{
		Clients: make(map[string]chan string),
	}
}

func (s *ClientState) AddClient(id string) chan string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Clients[id] = make(chan string)
	return s.Clients[id]
}

func (s *ClientState) RemoveClient(id string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Clients, id)
}

func (s *ClientState) GetAllClients() Clients {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.Clients
}
