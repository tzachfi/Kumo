package prompthub

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/tzachfi/kumo/server/internal/domain"
	"github.com/tzachfi/kumo/server/internal/prompthub/provider"
)

// promptFS embeds the prompt templates into the binary so the Hub has no runtime
// dependency on the filesystem layout.
//
//go:embed prompts/*.tmpl
var promptFS embed.FS

// Hub is the concrete MentorPrompter. It owns a Provider and the parsed prompt
// templates, and orchestrates the templating -> call -> parse lifecycle.
type Hub struct {
	provider  provider.Provider
	templates *template.Template
}

// NewHub parses the embedded templates and returns a ready Hub. It fails fast if
// any template is malformed.
func NewHub(p provider.Provider) (*Hub, error) {
	tmpls, err := template.ParseFS(promptFS, "prompts/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("prompthub: parse templates: %w", err)
	}
	return &Hub{provider: p, templates: tmpls}, nil
}

// render injects data into the named template and returns the compiled prompt.
func (h *Hub) render(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("prompthub: render %q: %w", name, err)
	}
	return buf.String(), nil
}

// GenerateJourney compiles the create_journey template from the given context,
// calls the provider, and parses the structured response into a domain.Journey.
func (h *Hub) GenerateJourney(ctx context.Context, req domain.JourneyContext) (*domain.Journey, error) {
	prompt, err := h.render("create_journey.tmpl", req)
	if err != nil {
		return nil, err
	}
	return executeStructured[domain.Journey](ctx, h.provider, prompt)
}

// Compile-time assurance that Hub satisfies the MentorPrompter contract.
var _ MentorPrompter = (*Hub)(nil)
