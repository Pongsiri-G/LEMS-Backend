package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod	= (pongWait * 9) / 10
)

// Client represents one connected WebSocket client.
type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	userID string
}

func newClient(conn *websocket.Conn, hub *Hub, userID string) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  hub,
		userID: userID,
	}
}

// ReadPump listens for messages from the client. 
func (c *Client) ReadPump() error { 
	defer func() { 
		c.hub.Unregister(c) 
		c.conn.Close() 
	}() 
	
	for { 
		if _, _, err := c.conn.NextReader(); err != nil { 
			return err 
		} 
	} 
} 

// WritePump sends broadcast messages to the client. 
func (c *Client) WritePump() error { 
	defer func () { 
		c.conn.Close() 
	}() 
	
	for msg := range c.send { 
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil { 
			return err 
		} 
	} 
	
	return nil 
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, userID string) error { 
	upgrader := websocket.Upgrader{ 
		CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins 
	} 
	
	conn, err := upgrader.Upgrade(w, r, nil) 
	
	if err != nil { 
		log.Error().Msgf("upgrade: %s", err.Error()) 
		return err 
	} 
	
	client := newClient(conn, hub, userID) 
	hub.Register(client) 

	go func() {
		if err := client.ReadPump(); err != nil {
			log.Error().Msgf("read pump ended: %s", err.Error())
		}
	}()
	go func() {
		if err := client.WritePump(); err != nil {
			log.Error().Msgf("write pump ended: %s", err.Error())
		}
	}()
	
	return nil
}
