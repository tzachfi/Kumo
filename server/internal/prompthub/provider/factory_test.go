package provider

import (
	"strings"
	"testing"

	"github.com/tzachfi/kumo/server/internal/prompthub/secrets"
)

func TestNewFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr string
		want    string // "mock", "llm", or ""
	}{
		{
			name: "default mock when provider unset",
			env: map[string]string{
				EnvProvider: "",
			},
			want: "mock",
		},
		{
			name: "explicit mock",
			env: map[string]string{
				EnvProvider: "mock",
			},
			want: "mock",
		},
		{
			name: "llm with all config",
			env: map[string]string{
				EnvProvider:        "llm",
				EnvLLMBaseURL:      "https://api.openai.com/v1",
				EnvLLMModel:        "gpt-4o-mini",
				secrets.EnvAPIKey:  "sk-test",
			},
			want: "llm",
		},
		{
			name: "unknown provider",
			env: map[string]string{
				EnvProvider: "bogus",
			},
			wantErr: "unknown provider",
		},
		{
			name: "llm missing base url",
			env: map[string]string{
				EnvProvider:       "llm",
				EnvLLMModel:       "gpt-4o-mini",
				secrets.EnvAPIKey: "sk-test",
			},
			wantErr: EnvLLMBaseURL,
		},
		{
			name: "llm missing model",
			env: map[string]string{
				EnvProvider:        "llm",
				EnvLLMBaseURL:      "https://api.openai.com/v1",
				secrets.EnvAPIKey:  "sk-test",
			},
			wantErr: EnvLLMModel,
		},
		{
			name: "llm missing api key",
			env: map[string]string{
				EnvProvider:   "llm",
				EnvLLMBaseURL: "https://api.openai.com/v1",
				EnvLLMModel:   "gpt-4o-mini",
				secrets.EnvAPIKey: "",
			},
			wantErr: secrets.EnvAPIKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(EnvProvider, "")
			t.Setenv(EnvLLMBaseURL, "")
			t.Setenv(EnvLLMModel, "")
			t.Setenv(secrets.EnvAPIKey, "")

			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			p, err := NewFromEnv()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("NewFromEnv() = nil, want error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want substring %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewFromEnv() = %v, want nil", err)
			}

			switch tt.want {
			case "mock":
				if _, ok := p.(*MockProvider); !ok {
					t.Fatalf("NewFromEnv() = %T, want *MockProvider", p)
				}
			case "llm":
				if _, ok := p.(*OpenAICompatProvider); !ok {
					t.Fatalf("NewFromEnv() = %T, want *OpenAICompatProvider", p)
				}
			default:
				t.Fatalf("invalid test case: want %q", tt.want)
			}
		})
	}
}
