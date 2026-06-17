package journey

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/tzachfi/kumo/server/internal/domain"
)

func TestAssemble_assignsRootIDs(t *testing.T) {
	userID := uuid.New()
	req := domain.JourneyContext{UserID: userID, 
	}
	j := &domain.Journey{}

	got := Assemble(j, req)
	if got != j {
		t.Fatal("Assemble should return the same pointer")
	}
	if j.ID == uuid.Nil {
		t.Fatal("Journey ID is not assigned")
	}
	if j.UserID != userID {
		t.Errorf("UserID = %v, want %v", j.UserID, userID)
	}
	if j.CreatedAt.IsZero() {
		t.Fatal("CreatedAt is not assigned")
	}
	if time.Since(j.CreatedAt) > 5*time.Second {
		t.Errorf("CreatedAt = %v, expected within last 5s", j.CreatedAt)
	}
}

func TestAssemble_assignsNestedIDs(t *testing.T) {
	j := &domain.Journey{
		ID: uuid.New(),
		Milestones: []domain.Milestone{
			{
				Title: "Test Milestone",
				Tasks: []domain.Task{
					{Title: "Test Task"},
				},
			},
		},
	}
	req := domain.JourneyContext{}
	Assemble(j, req)
	for i := range j.Milestones {
		m := &j.Milestones[i]
		if m.ID == uuid.Nil {
			t.Fatal("Milestone ID is not assigned")
		}
		if m.JourneyID != j.ID {
			t.Errorf("Milestone JourneyID = %q, want %q", m.JourneyID, j.ID)
		}
		for k := range m.Tasks {
			task := &m.Tasks[k]
			if task.ID == uuid.Nil {
				t.Fatal("Task ID is not assigned")
			}
			if task.MilestoneID != m.ID {
				t.Errorf("Task MilestoneID = %q, want %q", task.MilestoneID, m.ID)
			}
		}
	}
}

func TestAssemble_preservesExistingIDs(t *testing.T) {
	journeyUUID := uuid.New()
	milestoneUUID := uuid.New()
	taskUUID := uuid.New()
	j := &domain.Journey{
		ID: journeyUUID,
		Milestones: []domain.Milestone{
			{
				ID: milestoneUUID,
				Tasks: []domain.Task{
					{ID: taskUUID},
				},
			},
		},
	}
	req := domain.JourneyContext{}

	got := Assemble(j, req)
	if got != j {
		t.Fatal("Assemble should return the same pointer")
	}
	if j.ID != journeyUUID {
		t.Errorf("Journey ID = %q, want %q", j.ID, journeyUUID)
	}
	for i := range j.Milestones {
		m := &j.Milestones[i]
		if m.ID != milestoneUUID {
			t.Errorf("Milestone ID = %q, want %q", m.ID, milestoneUUID)
		}
		for k := range m.Tasks {
			task := &m.Tasks[k]
			if task.ID != taskUUID {
				t.Errorf("Task ID = %q, want %q", task.ID, taskUUID)
			}
		}
	}
}