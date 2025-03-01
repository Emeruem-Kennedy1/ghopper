package services

import (
	"sync"
)

type ClientManager struct {
	mu      sync.Mutex
	clients sync.Map
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func (cm *ClientManager) StoreClient(userID string, client SpotifyClientInterface) {
	cm.clients.Store(userID, client)
}

func (cm *ClientManager) GetClient(userID string) (SpotifyClientInterface, bool) {
	client, exists := cm.clients.Load(userID)
	if !exists {
		return nil, false
	}

	return client.(SpotifyClientInterface), true
}

func (cm *ClientManager) DeleteClient(userID string) {
	cm.clients.Delete(userID)
}

func (cm *ClientManager) RemoveClient(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients.Delete(userID)
}
