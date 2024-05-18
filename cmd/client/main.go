package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/pik4ez/faraway-assignment/internal/config"
	"github.com/pik4ez/faraway-assignment/internal/pow"
)

// configPath is a path to the application config file.
const configPath = "config.json"

func main() {
	// Connect the server.
	addr := os.Getenv("ADDR")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Error connecting to server: %+v\n", err)
		return
	}
	defer conn.Close()

	// Read the app config.
	cfg, err := config.NewFromFile(configPath)
	if err != nil {
		fmt.Printf("failed to read config: %+v\n", err)
	}

	// Complete the PoW challenge.
	err = pow.CompletePowChallenge(conn, cfg.ChallengeStrLen)
	if err != nil {
		fmt.Printf("Failed to pass PoW challenge: %+v\n", err)
	}

	// Read the response from the server.
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}
	response := string(buf[:n])
	if strings.HasPrefix(response, "QUOTE") {
		fmt.Println("Received quote:", strings.TrimSpace(response[6:]))
	} else {
		fmt.Println("Unexpected response from server:", response)
	}
}
