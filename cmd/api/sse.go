package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Client chan []byte // now sends JSON []byte

type Broker struct {
	clients map[Client]bool
	lock    sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[Client]bool),
	}
}

func (b *Broker) AddClient(c Client) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.clients[c] = true
}

func (b *Broker) RemoveClient(c Client) {
	b.lock.Lock()
	defer b.lock.Unlock()
	delete(b.clients, c)
	close(c)
}

func (b *Broker) NotifyAll(data interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	for c := range b.clients {
		select {
		case c <- jsonData:
		default:
			// Drop client if not reading
			delete(b.clients, c)
			close(c)
		}
	}
}

func (b *Broker) SSEHandler(w http.ResponseWriter, r *http.Request) {
	//token := r.URL.Query().Get("token")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	client := make(Client, 10)
	b.AddClient(client)
	defer b.RemoveClient(client)

	ctx := r.Context()

	// Heartbeat goroutine
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Comment event to keep connection alive
				fmt.Fprintf(w, ": ping\n\n")
				flusher.Flush()
			}
		}
	}()

	for {
		select {
		case msg := <-client:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}

// Example endpoint: notify with JSON
func (b *Broker) NotifyHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("msg")
	if message == "" {
		http.Error(w, "Missing msg param", http.StatusBadRequest)
		return
	}

	// Example structured data
	event := map[string]interface{}{
		"event":   "update",
		"message": message,
		"time":    time.Now().Format(time.RFC3339),
	}

	b.NotifyAll(event)
	fmt.Fprintln(w, "Message sent")
}

/*func main() {
	broker := NewBroker()
	http.HandleFunc("/events", broker.SSEHandler)
	http.HandleFunc("/notify", broker.NotifyHandler)

	fmt.Println("Server running at :8080")
	http.ListenAndServe(":8080", nil)
}
*/
