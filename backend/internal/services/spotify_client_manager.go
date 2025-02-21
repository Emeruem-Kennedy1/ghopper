package services

import (
	"sync"

	"github.com/zmb3/spotify"
)

type ClientManager struct {
	mu      sync.Mutex
	clients sync.Map
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func (cm *ClientManager) StoreClient(userID string, client *spotify.Client) {
	cm.clients.Store(userID, client)
}

func (cm *ClientManager) GetClient(userID string) (*spotify.Client, bool) {
	cleint, exists := cm.clients.Load(userID)
	if !exists {
		return nil, false
	}

	return cleint.(*spotify.Client), true
}

func (cm *ClientManager) DeleteClient(userID string) {
	cm.clients.Delete(userID)
}

func (cm *ClientManager) RemoveClient(userID string) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.clients.Delete(userID)
}