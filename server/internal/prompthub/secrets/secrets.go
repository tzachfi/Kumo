// Package secrets abstracts how the Prompt Hub loads sensitive values such as
// LLM API keys. Phase 2 reads from environment variables; a vault backend can
// implement the same Reader interface later without changing callers.
package secrets

// Reader loads runtime secrets for LLM providers.
type Reader interface {
	APIKey() (string, error)
}
