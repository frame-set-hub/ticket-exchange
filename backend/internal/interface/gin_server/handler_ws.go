package gin_server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TicketX/backend/internal/use_case"
	"github.com/TicketX/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// HandleWebSocket handles WebSocket connections for chat
func (s *GinServer) HandleWebSocket(c *gin.Context) {
	// 1. Get transaction ID
	txIDStr := c.Param("transaction_id")
	txID, err := strconv.ParseUint(txIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	// 2. Authenticate via query param token
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 3. Check authorization — user must be buyer, seller, or admin
	txResult, err := s.useCase.GetTransactionByID(c.Request.Context(), uint(txID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	tx := txResult.Transaction
	if claims.UserID != tx.BuyerID && claims.UserID != tx.SellerID && claims.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 4. Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	roomID := fmt.Sprintf("tx:%d", txID)
	client := &Client{
		hub:    s.hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		roomID: roomID,
		userID: claims.UserID,
	}

	s.hub.register <- client

	// Capture values for goroutine (don't use c.Request.Context() — it's cancelled after handler returns)
	senderUsername := claims.Username
	useCase := s.useCase
	hub := s.hub

	// 5. Start pumps
	go client.writePump()
	go client.readPump(func(cl *Client, msg *IncomingMessage) {
		// Determine receiver
		receiverID := tx.BuyerID
		if cl.userID == tx.BuyerID {
			receiverID = tx.SellerID
		}

		// Persist message — use background context since HTTP context is already done
		result, err := useCase.CreateMessage(context.Background(), &use_case.CreateMessageParams{
			TransactionID: uint(txID),
			SenderID:      cl.userID,
			ReceiverID:    receiverID,
			Content:       msg.Content,
			AttachmentURL: msg.AttachmentURL,
		})
		if err != nil {
			return
		}

		// Add sender username
		result.Message.SenderUsername = senderUsername

		// Broadcast to room
		msgJSON, _ := json.Marshal(result.Message)
		hub.broadcast <- &BroadcastMessage{
			RoomID: roomID,
			Data:   msgJSON,
		}
	})
}
