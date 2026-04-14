package api

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins; restrict in production via config
	},
}

// Hub manages all WebSocket connections and broadcasts
type Hub struct {
	mu      sync.RWMutex
	clients map[*wsClient]bool
}

type wsClient struct {
	conn *websocket.Conn
	send chan []byte
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{clients: make(map[*wsClient]bool)}
}

// Broadcast sends a message to all connected clients
func (hub *Hub) Broadcast(msg []byte) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	for c := range hub.clients {
		select {
		case c.send <- msg:
		default:
			// slow client; drop message
		}
	}
}

func (hub *Hub) register(c *wsClient) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.clients[c] = true
}

func (hub *Hub) unregister(c *wsClient) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	delete(hub.clients, c)
	close(c.send)
}

// WebSocketHandler upgrades HTTP to WebSocket and streams live KPI/metric events
func (h *Handlers) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade: %v", err)
		return
	}

	client := &wsClient{conn: conn, send: make(chan []byte, 256)}
	h.hub.register(client)

	// writer goroutine
	go func() {
		defer func() {
			conn.Close()
			h.hub.unregister(client)
		}()
		for msg := range client.send {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	// reader goroutine (keep-alive ping/pong)
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
