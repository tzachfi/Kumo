package prompthub

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/tzachfi/kumo/server/internal/domain"
	"github.com/tzachfi/kumo/server/internal/prompthub/provider"
)

const validJourneyJSON = `{
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

func newTestContext() domain.JourneyContext {
	return domain.JourneyContext{
		UserID:      uuid.New(),
		Ambition:    "run a 10k under 60 minutes",
		Domain:      domain.Fitness,
		Deadline:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		Constraints: []string{"only train 3 days/week"},
	}
}

func TestGenerateJourney_Success(t *testing.T) {
	hub, err := NewHub(provider.NewMockProvider(validJourneyJSON))
	if err != nil {
		t.Fatalf("NewHub: %v", err)
	}

	journey, err := hub.GenerateJourney(context.Background(), newTestContext())
	if err != nil {
		t.Fatalf("GenerateJourney: %v", err)
	}

	if journey.Title != "10k under 60 mins" {
		t.Errorf("Title = %q, want %q", journey.Title, "10k under 60 mins")
	}
	if journey.Domain != domain.Fitness {
		t.Errorf("Domain = %q, want %q", journey.Domain, domain.Fitness)
	}
	if got := len(journey.Milestones); got != 1 {
		t.Fatalf("len(Milestones) = %d, want 1", got)
	}
	if got := len(journey.Milestones[0].Tasks); got != 1 {
		t.Fatalf("len(Tasks) = %d, want 1", got)
	}
	if journey.Milestones[0].Tasks[0].State != domain.TaskPending {
		t.Errorf("Task state = %q, want %q", journey.Milestones[0].Tasks[0].State, domain.TaskPending)
	}
}

func TestGenerateJourney_RetryOnceThenSucceeds(t *testing.T) {
	mock := provider.NewMockProvider("this is not json", validJourneyJSON)

	hub, err := NewHub(mock)
	if err != nil {
		t.Fatalf("NewHub: %v", err)
	}

	journey, err := hub.GenerateJourney(context.Background(), newTestContext())
	if err != nil {
		t.Fatalf("GenerateJourney: %v", err)
	}
	if journey.Title != "10k under 60 mins" {
		t.Errorf("Title = %q, want %q", journey.Title, "10k under 60 mins")
	}

	if mock.CallCount() != 2 {
		t.Errorf("provider called %d times, want exactly 2 (initial + one retry)", mock.CallCount())
	}

	// The retry prompt must feed the parse failure back to the model.
	retryPrompt := mock.Prompts[1]
	if !strings.Contains(retryPrompt, "could not be parsed as JSON") {
		t.Errorf("retry prompt did not include the parse error feedback: %q", retryPrompt)
	}
}

func TestGenerateJourney_BadTwiceReturnsWrappedError(t *testing.T) {
	mock := provider.NewMockProvider("nope", "still not json")

	hub, err := NewHub(mock)
	if err != nil {
		t.Fatalf("NewHub: %v", err)
	}

	journey, err := hub.GenerateJourney(context.Background(), newTestContext())
	if err == nil {
		t.Fatalf("expected error, got journey: %+v", journey)
	}
	if journey != nil {
		t.Errorf("expected nil journey on error, got %+v", journey)
	}

	// Must stop after exactly one retry (no aggressive retry loop).
	if mock.CallCount() != 2 {
		t.Errorf("provider called %d times, want exactly 2", mock.CallCount())
	}

	// Error must carry retry context and remain unwrappable to the underlying cause.
	if !strings.Contains(err.Error(), "after one retry") {
		t.Errorf("error %q missing retry context", err.Error())
	}
	if errors.Unwrap(err) == nil {
		t.Errorf("expected a wrapped underlying error, got none")
	}
}
