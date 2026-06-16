package domain

// AmbitionDomain defines the category of the goal. The application stays
// agnostic to what each domain means; it only steps through the tree.
type AmbitionDomain string

const (
	Fitness   AmbitionDomain = "FITNESS"
	Technical AmbitionDomain = "TECHNICAL"
	Creative  AmbitionDomain = "CREATIVE"
)

// Valid reports whether the value is a recognized AmbitionDomain. String-typed
// enums are not constrained at compile time, so callers validate at boundaries.
func (d AmbitionDomain) Valid() bool {
	switch d {
	case Fitness, Technical, Creative:
		return true
	default:
		return false
	}
}

// JourneyState represents where the user stands in their overall timeline.
type JourneyState string

const (
	Initializing JourneyState = "INITIALIZING"
	Active       JourneyState = "ACTIVE"
	Paused       JourneyState = "PAUSED"
	Completed    JourneyState = "COMPLETED"
)

func (s JourneyState) Valid() bool {
	switch s {
	case Initializing, Active, Paused, Completed:
		return true
	default:
		return false
	}
}

// MilestoneState gates progression through the journey's phases.
type MilestoneState string

const (
	MilestoneLocked   MilestoneState = "LOCKED"
	MilestoneActive   MilestoneState = "ACTIVE"
	MilestoneFinished MilestoneState = "FINISHED"
)

func (s MilestoneState) Valid() bool {
	switch s {
	case MilestoneLocked, MilestoneActive, MilestoneFinished:
		return true
	default:
		return false
	}
}

// TaskState is the execution-layer state machine for individual actionable items.
type TaskState string

const (
	TaskPending   TaskState = "PENDING"
	TaskCompleted TaskState = "COMPLETED"
	TaskMissed    TaskState = "MISSED"
	TaskSkipped   TaskState = "SKIPPED"
)

func (s TaskState) Valid() bool {
	switch s {
	case TaskPending, TaskCompleted, TaskMissed, TaskSkipped:
		return true
	default:
		return false
	}
}
