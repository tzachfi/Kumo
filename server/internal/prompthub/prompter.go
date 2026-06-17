// Package prompthub is the core execution engine for interacting with external
// LLM APIs. It is domain-aware but business-agnostic: it handles data contracts,
// template compilation, API communication, and parsing, but never business rules
// such as authorization or whether a goal is realistic.
package prompthub

import (
	"context"

	"github.com/tzachfi/kumo/server/internal/domain"
)

// MentorPrompter exposes explicit, strongly typed generation methods grouped by
// core domain entities. Keeping methods typed (rather than one generic
// interface{} call) preserves Go's compile-time safety for callers.
type MentorPrompter interface {
	GenerateJourney(ctx context.Context, req domain.JourneyContext) (*domain.Journey, error)
	// GenerateAvatar(ctx context.Context, req domain.AvatarContext) (*domain.AvatarProfile, error) // later phase
}
