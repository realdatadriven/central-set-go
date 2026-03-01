package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Client is a channel that receives JSON bytes to send to a browser
type Client chan []byte

type Broker struct {
	clients map[Client]bool
	mu      sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[Client]bool),
	}
}

func (b *Broker) AddClient(c Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[c] = true
}

func (b *Broker) RemoveClient(c Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[c]; ok {
		delete(b.clients, c)
		close(c)
	}
}

// NotifyAll broadcasts the same message to all connected clients
func (b *Broker) NotifyAll(data any) {
	b.mu.Lock()
	defer b.mu.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return
	}

	for c := range b.clients {
		select {
		case c <- jsonData:
			// sent
		default:
			// client is slow or dead → drop it
			delete(b.clients, c)
			close(c)
		}
	}
}

// SSEHandler serves the /events stream
func (b *Broker) SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // adjust for production
	/*token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	fmt.Println(token)*/
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	client := make(Client, 10) // buffered → helps when browser is a bit slow
	b.AddClient(client)
	defer b.RemoveClient(client)
	// Heartbeat goroutine (keeps nginx / proxies happy)
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(w, ": ping\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}()
	// Main send loop
	for {
		select {
		case msg, ok := <-client:
			if !ok {
				return // channel closed → we're already removed
			}
			// msg is already []byte from json.Marshal
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

// NotifyHandler – example way to trigger a broadcast (GET /notify?msg=hello)
func (b *Broker) NotifyHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		http.Error(w, "missing ?msg= parameter", http.StatusBadRequest)
		return
	}

	event := map[string]any{
		"event":   "message",
		"content": msg,
		"time":    time.Now().UTC().Format(time.RFC3339),
		"sender":  "server",
	}

	b.NotifyAll(event)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Message broadcasted to all clients")
}

func (app *application) NotifyHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("msg")
	if message == "" {
		http.Error(w, "Missing msg param", http.StatusBadRequest)
		return
	}
	fmt.Println("SSE NotifyHandler Test:", message)
	http.Error(w, message, http.StatusBadRequest)
	return
}

/*func main() {
	broker := NewBroker()
	http.HandleFunc("/events", broker.SSEHandler)
	http.HandleFunc("/notify", broker.NotifyHandler)

	fmt.Println("Server running at :8080")
	http.ListenAndServe(":8080", nil)
}
*/
