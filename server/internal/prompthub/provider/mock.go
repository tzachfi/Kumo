package provider

import "context"

// MockProvider is a test double that returns a predefined sequence of responses,
// one per call. It lets tests drive the Hub's lifecycle (including the
// retry-once-on-parse-error path) without any network access.
type MockProvider struct {
	// Responses are returned in order, one per Complete call. The error in each
	// entry, if set, is returned alongside an empty string.
	Responses []MockResponse

	// Prompts records every prompt passed to Complete, in order, so tests can
	// assert that the retry attempt appended the parse error to the prompt.
	Prompts []string

	calls int
}

// MockResponse is a single canned reply for one Complete call.
type MockResponse struct {
	Text string
	Err  error
}

// NewMockProvider builds a MockProvider that returns the given texts in order,
// each with a nil error. Use the struct literal directly to inject errors.
func NewMockProvider(texts ...string) *MockProvider {
	responses := make([]MockResponse, len(texts))
	for i, t := range texts {
		responses[i] = MockResponse{Text: t}
	}
	return &MockProvider{Responses: responses}
}

// Complete returns the next queued response. If the queue is exhausted it
// repeats the final response, which keeps simple single-call tests concise.
func (m *MockProvider) Complete(_ context.Context, prompt string) (string, error) {
	m.Prompts = append(m.Prompts, prompt)

	if len(m.Responses) == 0 {
		return "", nil
	}

	idx := m.calls
	if idx >= len(m.Responses) {
		idx = len(m.Responses) - 1
	}
	m.calls++

	resp := m.Responses[idx]
	if resp.Err != nil {
		return "", resp.Err
	}
	return resp.Text, nil
}

// CallCount reports how many times Complete was invoked.
func (m *MockProvider) CallCount() int {
	return m.calls
}
