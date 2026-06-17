package secrets

import (
	"fmt"
	"os"
)

// EnvAPIKey is the environment variable name for the LLM provider API key.
const EnvAPIKey = "KUMO_LLM_API_KEY"

// Env reads secrets from process environment variables.
type Env struct{}

// NewEnv returns an Env-backed Reader.
func NewEnv() *Env {
	return &Env{}
}

// APIKey returns KUMO_LLM_API_KEY or an error if it is unset or empty.
func (e *Env) APIKey() (string, error) {
	key := os.Getenv(EnvAPIKey)
	if key == "" {
		return "", fmt.Errorf("secrets: %s is not set", EnvAPIKey)
	}
	return key, nil
}
