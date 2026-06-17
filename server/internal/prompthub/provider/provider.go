// Package provider abstracts the third-party LLM infrastructure behind a single
// interface so the Prompt Hub depends only on a contract, never on a concrete
// vendor (Gemini, DeepSeek, etc.). Tiered routing and real HTTP clients are
// added in later phases.
package provider

import "context"

// Provider executes a compiled prompt against an LLM and returns the raw text
// completion. Parsing into structured Go types is the Hub's responsibility.
type Provider interface {
	Complete(ctx context.Context, prompt string) (string, error)
}
