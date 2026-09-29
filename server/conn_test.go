package server

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-zoox/websocket/conn"
	"github.com/gorilla/websocket"
)

// TestCloseByServerIsNotAnError tests that a connection closed by the server
// itself, e.g. closed by heartbeat timeout, is treated as a normal close,
// instead of an error event.
func TestCloseByServerIsNotAnError(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	var mu sync.Mutex
	var gotError error

	closedCh := make(chan struct{})

	s.OnError(func(conn conn.Conn, err error) error {
		mu.Lock()
		defer mu.Unlock()

		gotError = err
		return nil
	})
	s.OnClose(func(conn conn.Conn, code int, message string) error {
		close(closedCh)
		return nil
	})
	s.OnConnect(func(conn conn.Conn) error {
		// the server closes the connection by itself
		return conn.Close()
	})

	httpServer := httptest.NewServer(s)
	defer httpServer.Close()

	client, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(httpServer.URL, "http"), nil)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer client.Close()

	select {
	case <-closedCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for close event")
	}

	// events are handled asynchronously, wait for a while to make sure
	// that an error event, if any, is handled
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if gotError != nil {
		t.Fatalf("expected no error event, got: %v", gotError)
	}
}
