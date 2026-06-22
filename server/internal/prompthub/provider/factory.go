package provider

import (
	"fmt"
	"os"

	"github.com/tzachfi/kumo/server/internal/prompthub/secrets"
)

const (
	EnvProvider   = "KUMO_PROVIDER"
	EnvLLMBaseURL = "KUMO_LLM_BASE_URL"
	EnvLLMModel   = "KUMO_LLM_MODEL"
)

// mockSampleJourneyJSON is canned LLM output for offline mock mode (no network).
const mockSampleJourneyJSON = `{
  "title": "10k under 60 mins",
  "domain": "FITNESS",
  "state": "INITIALIZING",
  "deadline": "2026-09-01T00:00:00Z",
  "config": {"avatar_tier": 1, "theme_palette": "neon_runner"},
  "milestones": [
    {
      "title": "Endurance Base Building",
      "order": 0,
      "state": "ACTIVE",
      "tasks": [
        {
          "title": "5km Interval Run",
          "scheduled_at": "2026-07-01T07:00:00Z",
          "state": "PENDING",
          "details": {"target_pace": "5:45", "distance_meters": 5000}
        }
      ]
    }
  ]
}`

// NewFromEnv selects and constructs a Provider from process environment.
// KUMO_PROVIDER is "mock" (default) or "llm".
func NewFromEnv() (Provider, error) {
	kind := os.Getenv(EnvProvider)
	if kind == "" {
		kind = "mock"
	}

	switch kind {
	case "mock":
		return NewMockProvider(mockSampleJourneyJSON), nil
	case "llm":
		baseURL := os.Getenv(EnvLLMBaseURL)
		if baseURL == "" {
			return nil, fmt.Errorf("provider: %s is required when KUMO_PROVIDER=llm", EnvLLMBaseURL)
		}
		model := os.Getenv(EnvLLMModel)
		if model == "" {
			return nil, fmt.Errorf("provider: %s is required when KUMO_PROVIDER=llm", EnvLLMModel)
		}
		apiKey, err := secrets.NewEnv().APIKey()
		if err != nil {
			return nil, err
		}
		return NewOpenAICompatProvider(OpenAICompatConfig{
			BaseURL: baseURL,
			Model:   model,
			APIKey:  apiKey,
		})
	default:
		return nil, fmt.Errorf("provider: unknown provider %q", kind)
	}
}
