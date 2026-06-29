// Package store is the database boundary. It returns pure domain records with
// no SDUI or HTTP dependencies.
package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tzachfi/kumo/server/internal/domain"
)

// ErrNotFound is returned when a journey ID does not exist.
var ErrNotFound = errors.New("store: journey not found")
var ErrBadInput = errors.New("store: bad input")

// MilestoneRecord is a single milestone row.
type MilestoneRecord struct {
	ID    string
	Title string
	Order int
}

// JourneyRecord is the persisted journey shape the SDUI layer reads.
type JourneyRecord struct {
	ID         string
	Topic      string
	Progress   int // 0–100
	SeedColor  string
	Milestones []MilestoneRecord
}

// JourneyStore loads and persists journeys.
type JourneyStore interface {
	GetJourneyByID(ctx context.Context, id string) (*JourneyRecord, error)

	// CreateJourney inserts a journey shell (typically state=INITIALIZING, no milestones).
	CreateJourney(ctx context.Context, j *domain.Journey) error

	// SaveMilestones inserts milestones and updates journey state in one transaction.
	SaveMilestones(ctx context.Context, journeyID uuid.UUID, ms []domain.Milestone, state domain.JourneyState) error
}
