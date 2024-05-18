package middleware

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/pik4ez/faraway-assignment/internal/pow"
)

// activeConnections is the current number of active connections.
var activeConnections int32

// PowMiddleware handles the PoW challenge-response mechanism.
func PowMiddleware(next func(net.Conn, context.Context)) func(net.Conn, context.Context) {
	return func(conn net.Conn, ctx context.Context) {
		challengeStrLen, ok := ctx.Value("challengeStrLen").(int)
		if !ok {
			fmt.Fprintln(conn, "POW-FAILURE")
			conn.Close()
			return
		}
		difficulty, ok := ctx.Value("difficulty").(int)
		if !ok {
			fmt.Fprintln(conn, "POW-FAILURE")
			conn.Close()
			return
		}
		if !pow.RunPowChallenge(conn, challengeStrLen, difficulty) {
			fmt.Fprintln(conn, "POW-FAILURE")
			conn.Close()
			return
		}

		next(conn, ctx)
	}
}

// GlobalRateLimiteMiddleware manages the connection count and sends SERVER-BUSY if the number of simultaneous connections hit.
func GlobalRateLimiteMiddleware(next func(net.Conn, context.Context)) func(net.Conn, context.Context) {
	return func(conn net.Conn, ctx context.Context) {
		maxConnections, ok := ctx.Value("maxConnections").(int32)
		if !ok {
			fmt.Fprintln(conn, "RATE-LIMITER-ERROR")
			conn.Close()
			return
		}
		baseDifficulty, ok := ctx.Value("difficulty").(int)
		if !ok {
			fmt.Fprintln(conn, "RATE-LIMITER-ERROR")
			conn.Close()
			return
		}
		connectionsBatchSize, ok := ctx.Value("connectionsBatchSize").(int)
		if !ok {
			fmt.Fprintln(conn, "RATE-LIMITER-ERROR")
			conn.Close()
			return
		}

		// Check if max connections limit is reached.
		if atomic.LoadInt32(&activeConnections) >= maxConnections {
			fmt.Fprintln(conn, "SERVER-BUSY")
			conn.Close()
			return
		}

		// Adjust the number of active connections.
		currentConnections := atomic.AddInt32(&activeConnections, 1)
		defer atomic.AddInt32(&activeConnections, -1)

		// Determine difficulty based on current number of connections.
		difficulty := determineDifficulty(baseDifficulty, connectionsBatchSize, currentConnections)
		fmt.Printf("Current Connections: %d, Difficulty: %d\n", currentConnections, difficulty)

		next(conn, context.WithValue(ctx, "difficulty", difficulty))
	}
}

// determineDifficulty assigns the difficulty based on a current number of concurrent connections.
//
// Difficulty will increase with every next batch of connections.
func determineDifficulty(baseDifficulty, batchSize int, currentConnections int32) int {
	additionalDifficulty := int(currentConnections / int32(batchSize))
	return baseDifficulty + additionalDifficulty
}
