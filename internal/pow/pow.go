package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

// RunPowChallenge runs the challenge over a TCP connection.
func RunPowChallenge(conn net.Conn, challengeStrLen, difficulty int) bool {
	challenge := generateRandomString(challengeStrLen)
	fmt.Fprintf(conn, "POW-CHALLENGE %d %s\n", difficulty, challenge)

	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return false
	}

	nonce := strings.TrimSpace(string(buf[:n]))
	return VerifyPoW(challenge, nonce, difficulty)
}

// CompletePowChallenge takes care of completing the PoW challenge received from a TCP connection.
func CompletePowChallenge(conn net.Conn, challengeStrLen int) error {
	// Read the PoW challenge from the server.
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("error reading from server: %w", err)
	}
	challengeMessage := string(buf[:n])
	if !strings.HasPrefix(challengeMessage, "POW-CHALLENGE") {
		return fmt.Errorf("unexpected response from server: %s", challengeMessage)
	}

	// Extract the challenge from the message.
	challengeParts := strings.Split(challengeMessage, " ")
	if len(challengeParts) != 3 {
		return fmt.Errorf("malformed challenge message: %s", challengeMessage)
	}
	difficulty, err := strconv.Atoi(strings.TrimSpace(challengeParts[1]))
	if err != nil {
		return fmt.Errorf("malformed difficulty: %w", err)
	}
	challenge := strings.TrimSpace(challengeParts[2])

	// Solve the PoW challenge.
	nonce := findSolution(challengeStrLen, challenge, difficulty)

	// Send the nonce back to the server.
	fmt.Fprintf(conn, nonce+"\n")

	return nil
}

// VerifyPoW validates the answer to the challenge.
func VerifyPoW(challenge, nonce string, difficulty int) bool {
	hash := sha256.Sum256([]byte(challenge + nonce))
	hashString := hex.EncodeToString(hash[:])
	return strings.HasPrefix(hashString, strings.Repeat("0", difficulty))
}

// findSolution finds a solution to the challenge.
func findSolution(nonceSize int, challenge string, difficulty int) string {
	nonce := generateRandomString(nonceSize)
	for {
		if VerifyPoW(challenge, nonce, difficulty) {
			return nonce
		}
		nonce = generateRandomString(nonceSize)
	}
}

// generateRandomString creates a random string with a set length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
