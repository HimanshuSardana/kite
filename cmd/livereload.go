package cmd

import "sync"

// reloadHub broadcasts reload events to all connected SSE clients.
// Each client gets its own buffered channel so a slow client never
// blocks the rebuild path or other clients.
type reloadHub struct {
	mu   sync.RWMutex
	subs map[chan struct{}]struct{}
}

func newReloadHub() *reloadHub {
	return &reloadHub{subs: make(map[chan struct{}]struct{})}
}

// subscribe registers a new client and returns its channel.
// The channel is buffered (size 1) so broadcast never blocks.
func (h *reloadHub) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subs[ch] = struct{}{}
	return ch
}

// unsubscribe removes a client channel.
func (h *reloadHub) unsubscribe(ch chan struct{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subs, ch)
}

// broadcast delivers a reload event to every subscriber without blocking.
// With zero subscribers it returns immediately and leaks nothing.
func (h *reloadHub) broadcast() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// count returns the number of current subscribers (for tests).
func (h *reloadHub) count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.subs)
}
