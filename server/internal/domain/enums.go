package domain

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
