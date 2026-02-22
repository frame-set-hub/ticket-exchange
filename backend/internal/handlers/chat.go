package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/TicketX/backend/internal/database"
	"github.com/TicketX/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for POC
	},
}

type Client struct {
	Conn   *websocket.Conn
	UserID uint
	Role   string
}

// Hub manages active connections mapping TransactionID -> Room -> Clients
// For POC, simplify to global tracking
type ChatHub struct {
	mu      sync.Mutex
	Clients map[*Client]bool
}

var Hub = &ChatHub{
	Clients: make(map[*Client]bool),
}

func (h *ChatHub) BroadcastMessage(msg models.Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.Clients {
		if client.Role == string(models.RoleAdmin) || client.UserID == msg.SenderID || client.UserID == msg.ReceiverID {
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				client.Conn.Close()
				delete(h.Clients, client)
			}
		}
	}
}

func ConnectChat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	role, _ := c.Get("role")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		Conn:   conn,
		UserID: uint(userID.(float64)),
		Role:   role.(string),
	}

	Hub.mu.Lock()
	Hub.Clients[client] = true
	Hub.mu.Unlock()

	defer func() {
		Hub.mu.Lock()
		delete(Hub.Clients, client)
		Hub.mu.Unlock()
		conn.Close()
	}()

	for {
		var incomingMsg struct {
			TransactionID uint   `json:"transaction_id"`
			Content       string `json:"content"`
			AttachmentURL string `json:"attachment_url"`
			ReceiverID    uint   `json:"receiver_id"`
		}

		err := conn.ReadJSON(&incomingMsg)
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		// Save to db
		msg := models.Message{
			TransactionID: incomingMsg.TransactionID,
			SenderID:      client.UserID,
			ReceiverID:    incomingMsg.ReceiverID, // Admin handles who they reply to
			Content:       incomingMsg.Content,
			AttachmentURL: incomingMsg.AttachmentURL,
			CreatedAt:     time.Now(),
		}

		database.DB.Create(&msg)
		Hub.BroadcastMessage(msg)
	}
}
