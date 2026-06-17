package domain

import (
	"time"

	"github.com/google/uuid"
)

// JourneyContext is the strongly typed input the Prompt Hub injects into a
// prompt template to generate a Journey. It carries the raw ambition and
// constraints; the Hub is responsible for templating and parsing, not for
// judging whether the goal is realistic (that is business logic elsewhere).
type JourneyContext struct {
	UserID      uuid.UUID      `json:"user_id"`
	Ambition    string         `json:"ambition"` // free-text goal, e.g. "run a 10k under 60 minutes"
	Domain      AmbitionDomain `json:"domain"`
	Deadline    time.Time      `json:"deadline"`
	Constraints []string       `json:"constraints,omitempty"` // e.g. "only train 3 days/week", "no gym access"
}
