package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/tzachfi/kumo/server/internal/domain"
)

type mockEntry struct {
	journey    domain.Journey
	milestones []domain.Milestone
	progress   int
}

// MockStore is an in-memory JourneyStore for development and tests.
type MockStore struct {
	byID map[uuid.UUID]*mockEntry
}

// NewMock returns a ready mock store.
func NewMock() *MockStore {
	return &MockStore{
		byID: make(map[uuid.UUID]*mockEntry),
	}
}

// GetJourneyByID returns a journey from memory when it has been persisted.
func (s *MockStore) GetJourneyByID(ctx context.Context, id string) (*JourneyRecord, error) {
	_ = ctx

	journeyID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrNotFound
	}

	entry, ok := s.byID[journeyID]
	if !ok {
		return nil, ErrNotFound
	}

	rec := journeyToRecord(
		entry.journey.ID,
		entry.journey.Title,
		entry.progress,
		entry.journey.Config,
		entry.milestones,
	)
	return &rec, nil
}

// CreateJourney stores a journey shell in memory.
func (s *MockStore) CreateJourney(ctx context.Context, j *domain.Journey) error {
	_ = ctx

	if j == nil {
		return ErrBadInput
	}

	journey := *j
	s.byID[journey.ID] = &mockEntry{
		journey:  journey,
		progress: 0,
	}
	return nil
}

// SaveMilestones attaches milestones and updates journey state in memory.
func (s *MockStore) SaveMilestones(ctx context.Context, journeyID uuid.UUID, ms []domain.Milestone, state domain.JourneyState) error {
	_ = ctx

	entry, ok := s.byID[journeyID]
	if !ok {
		return ErrNotFound
	}

	entry.milestones = append([]domain.Milestone(nil), ms...)
	entry.journey.State = state
	entry.progress = 0
	return nil
}

var _ JourneyStore = (*MockStore)(nil)
