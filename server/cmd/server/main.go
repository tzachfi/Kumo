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

func main() {
	fmt.Println("kumo server starting (phase 2)")

	prov, err := provider.NewFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("using provider: %s\n", providerKind(prov))

	hub, err := prompthub.NewHub(prov)
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
	fmt.Println("generating journey (prompt → LLM → parse)...")
	raw, err := hub.GenerateJourney(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate journey: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("journey parsed: %q\n", raw.Title)

	// Step 2: Assembly — inject system-owned IDs and foreign keys.
	fmt.Println("assembling system fields (IDs, user_id, created_at)...")
	assembled := journey.Assemble(raw, req)

	// Step 3: Validation — reject invalid enum values before persistence.
	fmt.Println("validating enums...")
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
	fmt.Println("done")
}

func providerKind(p provider.Provider) string {
	switch p.(type) {
	case *provider.MockProvider:
		return "mock"
	case *provider.OpenAICompatProvider:
		return "llm"
	default:
		return "unknown"
	}
}
