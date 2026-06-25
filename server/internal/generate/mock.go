// Package generate holds POST-only generation mocks (theme + content delays
// are orchestrated in the HTTP handler).
package generate

import "github.com/tzachfi/kumo/server/internal/store"

// MockTheme returns seed color and hero image keyword for a topic.
func MockTheme(topic string) (seedColor, heroKeyword string) {
	return "#3d85c6", "running"
}

// MockMilestones returns generated milestones for a topic.
func MockMilestones(topic string) []store.MilestoneRecord {
	return []store.MilestoneRecord{
		{
			Title: "Base Building",
			Order: 1,
		},
		{
			Title: "Race Prep",
			Order: 2,
		},
	}
}
