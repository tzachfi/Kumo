package provider

import (
	"strings"
	"testing"
)

func TestNewOpenAICompatProvider_validation(t *testing.T) {
	valid := OpenAICompatConfig{
		BaseURL: "https://api.openai.com/v1",
		Model:   "gpt-4o-mini",
		APIKey:  "sk-test",
	}

	tests := []struct {
		name    string
		cfg     OpenAICompatConfig
		wantErr string
	}{
		{name: "valid config", cfg: valid, wantErr: ""},
		{
			name:    "missing BaseURL",
			cfg:     OpenAICompatConfig{Model: valid.Model, APIKey: valid.APIKey},
			wantErr: "BaseURL",
		},
		{
			name:    "missing Model",
			cfg:     OpenAICompatConfig{BaseURL: valid.BaseURL, APIKey: valid.APIKey},
			wantErr: "Model",
		},
		{
			name:    "missing APIKey",
			cfg:     OpenAICompatConfig{BaseURL: valid.BaseURL, Model: valid.Model},
			wantErr: "APIKey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewOpenAICompatProvider(tt.cfg)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("NewOpenAICompatProvider() = %v, want nil", err)
				}
				if p == nil {
					t.Fatal("NewOpenAICompatProvider() = nil, want provider")
				}
				return
			}
			if err == nil {
				t.Fatal("NewOpenAICompatProvider() = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestStripCodeFences(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "json fence",
			input: "```json\n{\"a\":1}\n```",
			want:  `{"a":1}`,
		},
		{
			name:  "plain json",
			input: `{"a":1}`,
			want:  `{"a":1}`,
		},
		{
			name:  "plain json with outer whitespace",
			input: "  {\"a\":1}  \n",
			want:  `{"a":1}`,
		},
		{
			name:  "fence without language tag",
			input: "```\n{\"a\":1}\n```",
			want:  `{"a":1}`,
		},
		{
			name:  "opening fence without closing",
			input: "```json\n{\"a\":1}",
			want:  "```json\n{\"a\":1}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripCodeFences(tt.input); got != tt.want {
				t.Errorf("stripCodeFences() = %q, want %q", got, tt.want)
			}
		})
	}
}
