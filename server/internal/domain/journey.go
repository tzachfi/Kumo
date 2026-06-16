package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Journey is the root contract representing the user's entire ambition lifecycle.
// Strict relational fields drive the state machine and progress tracking, while
// Config holds a loose, domain-specific jsonb payload the backend never inspects.
type Journey struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	UserID     uuid.UUID       `json:"user_id" db:"user_id"`
	Title      string          `json:"title" db:"title"` // e.g. "10k under 60 mins"
	Domain     AmbitionDomain  `json:"domain" db:"domain"`
	State      JourneyState    `json:"state" db:"state"`
	Deadline   time.Time       `json:"deadline" db:"deadline"`
	Config     json.RawMessage `json:"config,omitempty" db:"config"` // Polymorphic domain setup (pace, theme, avatar tier)
	Milestones []Milestone     `json:"milestones,omitempty"`         // Ordered timeline phases
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

// Milestone represents a phase of the journey. Milestones provide structural
// scaffolding so an LLM can recalculate one isolated block at a time rather than
// holding a flat list of every task in memory.
type Milestone struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	JourneyID uuid.UUID      `json:"journey_id" db:"journey_id"`
	Title     string         `json:"title" db:"title"` // e.g. "Endurance Base Building"
	Order     int            `json:"order" db:"sequence_order"`
	State     MilestoneState `json:"state" db:"state"`
	Tasks     []Task         `json:"tasks,omitempty"`
}

// Task is an individual actionable item on the schedule and the execution layer
// of the state machine. Details and ProofOfWork are domain-specific jsonb blobs
// consumed by the SDUI payload builder, not by core business logic.
type Task struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	MilestoneID uuid.UUID       `json:"milestone_id" db:"milestone_id"`
	Title       string          `json:"title" db:"title"` // e.g. "5km Interval Run"
	ScheduledAt time.Time       `json:"scheduled_at" db:"scheduled_at"`
	State       TaskState       `json:"state" db:"state"`
	Details     json.RawMessage `json:"details,omitempty" db:"details"`             // Domain payload (pace intervals, code docs)
	ProofOfWork json.RawMessage `json:"proof_of_work,omitempty" db:"proof_of_work"` // Validation data (Strava link, commit hash)
}
