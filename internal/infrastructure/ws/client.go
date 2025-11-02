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

	c.conn.SetReadLimit(512)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Print("unexpected websocket close:", err)
			}
			return err
		}
	}
} 

// WritePump sends broadcast messages to the client. 
func (c *Client) WritePump() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{}) // close handshake
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the channel
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}
			_, _ = w.Write(msg)
			if err := w.Close(); err != nil {
				return err
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
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
