// Package sdui defines the server-driven UI JSON contract and builders that
// translate store records into component trees at request time.
package sdui

// ComponentType is the React renderer discriminator.
type ComponentType string

const (
	TypeScreen         ComponentType = "Screen"
	TypeHero           ComponentType = "Hero"
	TypeSkeletonLoader ComponentType = "SkeletonLoader"
	TypeCard           ComponentType = "Card"
	TypeProgressBar    ComponentType = "ProgressBar"
)

// Component is a node in the SDUI tree returned to the client.
type Component struct {
	Type     ComponentType  `json:"type"`
	ID       string         `json:"id"`
	Props    map[string]any `json:"props,omitempty"`
	Children []Component    `json:"children,omitempty"`
}

// UpdatePayload is sent on SSE "update" events to patch the live tree.
type UpdatePayload struct {
	Action     string      `json:"action"`
	TargetID   string      `json:"targetId"`
	Components []Component `json:"components"`
}
