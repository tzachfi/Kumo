// Command server is a thin entrypoint for Kumo. It orchestrates the Phase 2
// pipeline: Prompt Hub (LLM parse) → journey assembly → validation → output.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/tzachfi/kumo/server/internal/domain"
	"github.com/tzachfi/kumo/server/internal/journey"
	"github.com/tzachfi/kumo/server/internal/prompthub"
	"github.com/tzachfi/kumo/server/internal/prompthub/provider"
)

const sampleJourneyJSON = `{
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

func main() {
	fmt.Println("kumo server starting (phase 2: hub → assemble → validate)")

	hub, err := prompthub.NewHub(provider.NewMockProvider(sampleJourneyJSON))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build prompt hub: %v\n", err)
		os.Exit(1)
	}

	req := domain.JourneyContext{
		UserID:      uuid.New(),
		Ambition:    "run a 10k under 60 minutes",
		Domain:      domain.Fitness,
		Deadline:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		Constraints: []string{"only train 3 days/week", "no gym access"},
	}

	ctx := context.Background()

	// Step 1: Prompt Hub — template, LLM call, JSON parse (LLM-shaped output).
	raw, err := hub.GenerateJourney(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate journey: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Assembly — inject system-owned IDs and foreign keys.
	assembled := journey.Assemble(raw, req)

	// Step 3: Validation — reject invalid enum values before persistence.
	if err := journey.Validate(assembled); err != nil {
		fmt.Fprintf(os.Stderr, "journey validation failed: %v\n", err)
		os.Exit(1)
	}

	pretty, err := json.MarshalIndent(assembled, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode journey: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("generated journey:\n%s\n", pretty)
}
