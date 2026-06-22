package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OpenAICompatConfig holds settings for an OpenAI-compatible chat/completions API.
// BaseURL and Model come from env so one client works with OpenAI, DeepSeek,
// Gemini (compat endpoint), and similar gateways.
type OpenAICompatConfig struct {
	BaseURL string       // e.g. https://api.openai.com/v1
	Model   string       // e.g. gpt-4o-mini
	APIKey  string       // Bearer token from secrets.Reader
	Client  *http.Client // optional; defaults to http.DefaultClient
}

// OpenAICompatProvider calls remote LLMs that expose POST /chat/completions
// with the standard OpenAI request/response shape. Vendors with different APIs
// need a separate Provider implementation.
type OpenAICompatProvider struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// NewOpenAICompatProvider validates cfg and returns a ready provider.
func NewOpenAICompatProvider(cfg OpenAICompatConfig) (*OpenAICompatProvider, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("provider: BaseURL is required")
	}
	if cfg.Model == "" {
		return nil, fmt.Errorf("provider: Model is required")
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("provider: APIKey is required")
	}

	client := cfg.Client
	if client == nil {
		client = http.DefaultClient
	}

	return &OpenAICompatProvider{
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		apiKey:  cfg.APIKey,
		client:  client,
	}, nil
}

// Complete sends prompt to the LLM and returns the raw completion text.
func (p *OpenAICompatProvider) Complete(ctx context.Context, prompt string) (string, error) {
	payload := chatRequest{
		Model:    p.model,
		Messages: []chatMessage{{Role: "user", Content: prompt}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("provider: marshal request: %w", err)
	}

	url := strings.TrimSuffix(p.baseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("provider: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("provider: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("provider: read body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("provider: status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var out chatResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("provider: unmarshal response: %w", err)
	}
	if len(out.Choices) == 0 || out.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("provider: empty choices or content")
	}

	return stripCodeFences(out.Choices[0].Message.Content), nil
}

// stripCodeFences removes markdown code fences from LLM output so JSON parsers see plain text.
func stripCodeFences(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "```") {
		return trimmed
	}

	rest := trimmed[3:]
	i := strings.Index(rest, "\n")
	if i < 0 {
		return trimmed
	}
	rest = rest[i+1:]

	close := strings.LastIndex(rest, "```")
	if close < 0 {
		return trimmed
	}

	return strings.TrimSpace(rest[:close])
}

var _ Provider = (*OpenAICompatProvider)(nil)
