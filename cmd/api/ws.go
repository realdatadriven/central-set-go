package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	activeConnections map[*websocket.Conn]bool
	mu                sync.Mutex
}

func (app *application) NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		activeConnections: make(map[*websocket.Conn]bool),
	}
}

func (m *ConnectionManager) Connect(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeConnections[conn] = true
}

func (m *ConnectionManager) Disconnect(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeConnections, conn)
	conn.Close()
}

func (m *ConnectionManager) SendPersonalMessage(conn *websocket.Conn, message map[string]interface{}) error {
	return conn.WriteJSON(message)
}

func (m *ConnectionManager) Broadcast(message map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for conn := range m.activeConnections {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Broadcast error: %v", err)
			m.Disconnect(conn)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // You can add origin checks here if necessary
	},
}

func (app *application) broadcastTableChange(manager *ConnectionManager, message map[string]interface{}) {
	manager.Broadcast(message)
}

func (app *application) websocketEndpoint(manager *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}
		defer manager.Disconnect(conn)

		manager.Connect(conn)

		for {
			var data map[string]interface{}
			err := conn.ReadJSON(&data)
			if err != nil {
				log.Printf("Read error: %v", err)
				break
			}
			manager.Broadcast(data)
		}
	}
}
