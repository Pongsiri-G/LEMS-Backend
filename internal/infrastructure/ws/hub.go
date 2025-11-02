package ws

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// Hub maintains active clients and broadcasts messages to all of them.
type Hub struct {
	clients map[*Client]bool
	mutex   sync.Mutex
	userClients map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		userClients: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[c] = true
	if _, ok := h.userClients[c.userID]; !ok {
		h.userClients[c.userID] = make(map[*Client]bool)
	}
	
	h.userClients[c.userID][c] = true

	log.Printf("user %s, total: %d", c.userID, len(h.clients))
}

func (h *Hub) Unregister(c *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.clients, c)
	if clients, ok := h.userClients[c.userID]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.userClients, c.userID)
		}
	}
		
	log.Print("❌ WebSocket disconnected:", len(h.clients), "clients")
}

func (h *Hub) Broadcast(msg []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	log.Info().Msg("Boardcasting to all client")

	for c := range h.clients {
		select {
		case c.send <- msg:
		default:
			delete(h.clients, c)
			close(c.send)
		}
	}
}

func (h *Hub) SendToUser(userID string, msg []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	log.Info().Str("user", userID).Msg("Boardcasting to user")

	if clients, ok:= h.userClients[userID]; ok {
		for c := range clients {
			select {
			case c.send <- msg:
			default:
				delete(h.clients, c)
				delete(h.userClients[userID], c)
				close(c.send)
			}
		}
	}	
}
