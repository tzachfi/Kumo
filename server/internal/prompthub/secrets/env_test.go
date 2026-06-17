package secrets

import (
	"strings"
	"testing"
)

func TestEnv_APIKey(t *testing.T) {
	t.Run("returns key when set", func(t *testing.T) {
		t.Setenv(EnvAPIKey, "sk-test-key")
		got, err := NewEnv().APIKey()
		if err != nil {
			t.Fatalf("APIKey() = %v, want nil", err)
		}
		if got != "sk-test-key" {
			t.Errorf("APIKey() = %q, want %q", got, "sk-test-key")
		}
	})

	t.Run("error when unset", func(t *testing.T) {
		t.Setenv(EnvAPIKey, "")
		_, err := NewEnv().APIKey()
		if err == nil {
			t.Fatal("APIKey() = nil, want error")
		}
		if !strings.Contains(err.Error(), EnvAPIKey) {
			t.Errorf("error = %q, want mention of %q", err.Error(), EnvAPIKey)
		}
	})
}

// Compile-time check that Env satisfies Reader.
var _ Reader = (*Env)(nil)
