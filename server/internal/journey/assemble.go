package journey

import (
	"time"

	"github.com/google/uuid"

	"github.com/tzachfi/kumo/server/internal/domain"
)

// Assemble injects system-owned fields into an LLM-parsed Journey.
// It mutates j in place and returns the same pointer.
func Assemble(j *domain.Journey, req domain.JourneyContext) *domain.Journey {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	j.UserID = req.UserID

	if j.CreatedAt.IsZero() {
		j.CreatedAt = time.Now().UTC()
	}
	if j.Deadline.IsZero() {
		j.Deadline = req.Deadline
	}

	for i := range j.Milestones {
		m := &j.Milestones[i]
		if m.ID == uuid.Nil {
			m.ID = uuid.New()
		}
		m.JourneyID = j.ID

		for k := range m.Tasks {
			t := &m.Tasks[k]
			if t.ID == uuid.Nil {
				t.ID = uuid.New()
			}
			t.MilestoneID = m.ID
		}
	}

	return j
}
