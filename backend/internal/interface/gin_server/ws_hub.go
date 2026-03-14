package gin_server

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// BroadcastMessage represents a message to broadcast to a room
type BroadcastMessage struct {
	RoomID string
	Data   []byte
}

// Client represents a single WebSocket connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	roomID string
	userID uint
}

// Hub manages WebSocket rooms and client connections
type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Run starts the hub event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.roomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.rooms, client.roomID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.rooms[msg.RoomID]
			for client := range clients {
				select {
				case client.send <- msg.Data:
				default:
					// Client buffer full — disconnect
					close(client.send)
					delete(clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// IncomingMessage represents a message received from WebSocket client
type IncomingMessage struct {
	Content       string `json:"content"`
	AttachmentURL string `json:"attachment_url"`
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump(onMessage func(client *Client, msg *IncomingMessage)) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			break
		}

		var msg IncomingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		if msg.Content == "" {
			continue
		}

		onMessage(c, &msg)
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
