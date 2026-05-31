package sse

import "sync"

type Hub struct {
	mu      sync.RWMutex
	clients map[string][]chan string
}

var Default = &Hub{clients: make(map[string][]chan string)}

func (h *Hub) Subscribe(userID string) chan string {
	ch := make(chan string, 4)
	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], ch)
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(userID string, ch chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for i, c := range h.clients[userID] {
		if c == ch {
			h.clients[userID] = append(h.clients[userID][:i], h.clients[userID][i+1:]...)
			break
		}
	}
	close(ch)
}

func (h *Hub) Notify(userID, event string) {
	h.mu.RLock()
	chans := make([]chan string, len(h.clients[userID]))
	copy(chans, h.clients[userID])
	h.mu.RUnlock()
	for _, ch := range chans {
		select {
		case ch <- event:
		default:
		}
	}
}
