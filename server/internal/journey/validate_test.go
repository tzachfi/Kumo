package journey

import (
	"strings"
	"testing"

	"github.com/tzachfi/kumo/server/internal/domain"
)

func validJourney() *domain.Journey {
	return &domain.Journey{
		Title: "10k under 60 mins",
		State: domain.Initializing,
		Milestones: []domain.Milestone{
			{
				Title: "Base Building",
				State: domain.MilestoneActive,
				Tasks: []domain.Task{
					{Title: "5km Run", State: domain.TaskPending},
				},
			},
		},
	}
}

func TestValidate_validJourney(t *testing.T) {
	if err := Validate(validJourney()); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestValidate_nilJourney(t *testing.T) {
	if err := Validate(nil); err == nil {
		t.Fatal("Validate(nil) = nil, want error")
	}
}

func TestValidate_invalidStates(t *testing.T) {
	tests := []struct {
		name    string
		journey *domain.Journey
		wantSub string
	}{
		{
			name: "invalid journey state",
			journey: func() *domain.Journey {
				j := validJourney()
				j.State = "BOGUS"
				return j
			}(),
			wantSub: "invalid state",
		},
		{
			name: "invalid milestone state",
			journey: func() *domain.Journey {
				j := validJourney()
				j.Milestones[0].State = "BOGUS"
				return j
			}(),
			wantSub: "milestone[0]: invalid state",
		},
		{
			name: "invalid task state",
			journey: func() *domain.Journey {
				j := validJourney()
				j.Milestones[0].Tasks[0].State = "BOGUS"
				return j
			}(),
			wantSub: "milestone[0].task[0]: invalid state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.journey)
			if err == nil {
				t.Fatal("Validate() = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.wantSub)
			}
		})
	}
}
