package store

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/tzachfi/kumo/server/internal/domain"
)

const defaultSeedColor = "#3d85c6"

// journeyToRecord maps normalized domain rows to the SDUI read projection.
func journeyToRecord(
	id uuid.UUID,
	title string,
	progressPct int,
	config json.RawMessage,
	milestones []domain.Milestone,
) JourneyRecord {
	var cfg struct {
		SeedColor string `json:"seed_color"`
	}

	seedColor := defaultSeedColor
	if err := json.Unmarshal(config, &cfg); err == nil && cfg.SeedColor != "" {
		seedColor = cfg.SeedColor
	}

	ms := make([]MilestoneRecord, 0, len(milestones))
	for _, m := range milestones {
		ms = append(ms, MilestoneRecord{
			ID:    m.ID.String(),
			Title: m.Title,
			Order: m.Order,
		})
	}

	return JourneyRecord{
		ID:         id.String(),
		Topic:      title,
		Progress:   progressPct,
		SeedColor:  seedColor,
		Milestones: ms,
	}
}
