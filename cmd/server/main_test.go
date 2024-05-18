package main

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/pik4ez/faraway-assignment/internal/config"
)

// TestGlobalRateLimiter checks whether server returns SERVER-BUSY when hitting the maximum number of connections.
func TestGlobalRateLimiter(t *testing.T) {
	addr := "127.0.0.1:12346"

	cfg := config.Config{
		MaxSimultaneousConnections: 2,
		ConnectionsBatchSize:       1,
		BaseDifficulty:             1,
		ChallengeStrLen:            3,
	}

	// Start server in a goroutine, give it time to start.
	go startServer(addr, cfg)

	time.Sleep(1 * time.Second)

	// Set number of connections higher than the server maximum.
	connections := 3

	rateLimitHits := 0
	for i := 0; i < connections; i++ {
		go func() {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Errorf("Error connecting to server: %v", err)
				return
			}
			defer conn.Close()

			// Read the response from the server
			buf := make([]byte, 256)
			n, err := conn.Read(buf)
			if err != nil {
				t.Errorf("Error reading from server: %v", err)
				return
			}
			response := strings.TrimSpace(string(buf[:n]))

			if response == "SERVER-BUSY" {
				rateLimitHits++
			} else if !strings.HasPrefix(response, "POW-CHALLENGE") {
				t.Errorf("Unexpected response from server: %s", response)
			}
		}()
	}

	// Give clients time to connect.
	time.Sleep(3 * time.Second)

	if rateLimitHits == 0 {
		t.Errorf("Expected some connections to be refused due to server busy, but got none")
	} else {
		t.Logf("Server responded with SERVER-BUSY to %d connections", rateLimitHits)
	}
}
