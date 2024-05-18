package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"

	"github.com/pik4ez/faraway-assignment/internal/config"
	"github.com/pik4ez/faraway-assignment/internal/middleware"
)

const configPath = "config.json"

// quotes is the source of wisdom.
var quotes = []string{
	"The only true wisdom is in knowing you know nothing. - Socrates",
	"The unexamined life is not worth living. - Socrates",
	"To be yourself in a world that is constantly trying to make you something else is the greatest accomplishment. - Ralph Waldo Emerson",
	"In the end, we will remember not the words of our enemies, but the silence of our friends. - Martin Luther King Jr.",
	"To live is the rarest thing in the world. Most people exist, that is all. - Oscar Wilde",
}

// quoteHandler responds with quotes.
func quoteHandler(conn net.Conn, ctx context.Context) {
	defer conn.Close()

	quote := quotes[rand.Intn(len(quotes))]
	fmt.Fprintf(conn, "QUOTE %s\n", quote)
}

// startServer runs a server.
func startServer(addr string, cfg config.Config) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx = context.WithValue(ctx, "challengeStrLen", cfg.ChallengeStrLen)
		ctx = context.WithValue(ctx, "difficulty", cfg.BaseDifficulty)
		ctx = context.WithValue(ctx, "maxConnections", cfg.MaxSimultaneousConnections)
		ctx = context.WithValue(ctx, "connectionsBatchSize", cfg.ConnectionsBatchSize)

		go middleware.GlobalRateLimiteMiddleware(middleware.PowMiddleware(quoteHandler))(conn, ctx)
	}
}

// main is the entrypoint.
func main() {
	addr := os.Getenv("ADDR")
	cfg, err := config.NewFromFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config: %+v\n", err)
		os.Exit(1)
	}
	startServer(addr, cfg)
}
