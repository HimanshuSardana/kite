package cmd

import (
	"runtime"
	"testing"
	"time"
)

func recv(t *testing.T, ch chan struct{}, name string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatalf("subscriber %s did not receive reload", name)
	}
}

func TestBroadcastReachesAllSubscribers(t *testing.T) {
	h := newReloadHub()
	a := h.subscribe()
	defer h.unsubscribe(a)
	b := h.subscribe()
	defer h.unsubscribe(b)

	h.broadcast()

	recv(t, a, "a")
	recv(t, b, "b")
}

func TestBroadcastWithZeroSubscribers(t *testing.T) {
	h := newReloadHub()

	before := runtime.NumGoroutine()
	for i := 0; i < 100; i++ {
		h.broadcast()
	}
	// Broadcasts are synchronous and non-blocking: no goroutines spawned.
	// Give the scheduler a moment, then assert no growth.
	time.Sleep(50 * time.Millisecond)
	after := runtime.NumGoroutine()
	if after > before {
		t.Fatalf("goroutine leak: before=%d after=%d", before, after)
	}
	if got := h.count(); got != 0 {
		t.Fatalf("expected 0 subscribers, got %d", got)
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	h := newReloadHub()
	a := h.subscribe()
	h.unsubscribe(a)

	if got := h.count(); got != 0 {
		t.Fatalf("expected 0 subscribers after unsubscribe, got %d", got)
	}

	h.broadcast()

	select {
	case <-a:
		t.Fatal("unsubscribed client received reload")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestBroadcastNonBlockingWhenSlowClient(t *testing.T) {
	h := newReloadHub()
	slow := h.subscribe()
	defer h.unsubscribe(slow)
	fast := h.subscribe()
	defer h.unsubscribe(fast)

	// Fill the slow client's buffer, then broadcast again: must not block.
	h.broadcast()
	h.broadcast()

	done := make(chan struct{})
	go func() {
		h.broadcast()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("broadcast blocked on slow client")
	}

	recv(t, fast, "fast")
}
