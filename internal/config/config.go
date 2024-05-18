package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config is the global configuration.
type Config struct {
	BaseDifficulty             int   `json:"baseDifficulty"`
	ChallengeStrLen            int   `json:"challengeStrLen"`
	ConnectionsBatchSize       int   `json:"connectionsBatchSize"`
	MaxSimultaneousConnections int32 `json:"maxSimultaneousConnections"`
}

// NewFromFile creates a new config.
func NewFromFile(path string) (Config, error) {
	var conf Config
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file %s", path)
	}
	json.Unmarshal(raw, &conf)
	return conf, nil
}
