package store

// MockStore is an in-memory JourneyStore for development and tests.
type MockStore struct{}

// NewMock returns a ready mock store.
func NewMock() *MockStore {
	return &MockStore{}
}

// GetJourneyByID returns fixed mock data for known IDs.
func (s *MockStore) GetJourneyByID(id string) (*JourneyRecord, error) {
	if id == "" || id == "missing" {
		return nil, ErrNotFound
	}

	return &JourneyRecord{
		ID:        id,
		Topic:     "10km Training",
		Progress:  45,
		SeedColor: "#3d85c6",
		Milestones: []MilestoneRecord{
			{
				ID:    "base-building",
				Title: "Base Building",
				Order: 1,
			},
			{
				ID:    "race-prep",
				Title: "Race Prep",
				Order: 2,
			},
		},
	}, nil
}

var _ JourneyStore = (*MockStore)(nil)
