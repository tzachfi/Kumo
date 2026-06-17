package prompthub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tzachfi/kumo/server/internal/prompthub/provider"
)

// executeStructured runs the shared Prompt Hub lifecycle for any target type T:
//
//  1. Send the compiled prompt to the provider.
//  2. Unmarshal the raw completion into T.
//  3. On a JSON parse failure only, retry exactly once, appending the parse
//     error to the prompt so the model can self-correct its syntax.
//
// Per the resiliency policy, transport/API errors are NOT retried (LLM calls are
// slow and expensive); only a parse failure earns a single self-correction
// attempt. All terminal errors are wrapped with %w for context-rich propagation.
func executeStructured[T any](ctx context.Context, p provider.Provider, prompt string) (*T, error) {
	raw, err := p.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("prompthub: llm completion failed: %w", err)
	}

	var out T
	parseErr := json.Unmarshal([]byte(raw), &out)
	if parseErr == nil {
		return &out, nil
	}

	// Smart parsing retry: feed the parse error back to the model, once.
	retryPrompt := fmt.Sprintf(
		"%s\n\nYour previous response could not be parsed as JSON (%s). "+
			"Respond again with ONLY valid JSON that matches the required schema.",
		prompt, parseErr,
	)

	raw, err = p.Complete(ctx, retryPrompt)
	if err != nil {
		return nil, fmt.Errorf("prompthub: llm completion failed on parse-retry: %w", err)
	}

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("prompthub: failed to parse llm response after one retry: %w", err)
	}

	return &out, nil
}
