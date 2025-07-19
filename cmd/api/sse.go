package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Client chan string

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

func (b *Broker) NotifyAll(msg string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for c := range b.clients {
		select {
		case c <- msg:
		default:
			// Drop client if it's not reading
			delete(b.clients, c)
			close(c)
		}
	}
}

func (b *Broker) SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	client := make(Client, 10) // buffered channel
	b.AddClient(client)
	defer b.RemoveClient(client)

	ctx := r.Context()

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

func (b *Broker) NotifyHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("msg")
	if message == "" {
		http.Error(w, "Missing msg param", http.StatusBadRequest)
		return
	}
	b.NotifyAll(message)
	fmt.Fprintln(w, "Message sent")
}
