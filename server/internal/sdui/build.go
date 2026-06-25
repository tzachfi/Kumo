package sdui

import (
	"fmt"

	"github.com/tzachfi/kumo/server/internal/store"
)

// BuildDashboardScreen builds the full GET /api/journey/{id} tree.
func BuildDashboardScreen(rec store.JourneyRecord) (Component, error) {
	vars, err := PaletteFromSeed(rec.SeedColor)
	if err != nil {
		return Component{}, err
	}

	screen := Component{Type: TypeScreen, ID: "dashboard", Props: map[string]any{"style": vars}}

	screen.Children = append(screen.Children, buildHero(rec.Topic, ""))
	screen.Children = append(screen.Children, buildProgressBar(rec.Progress))
	screen.Children = append(screen.Children, buildCards(rec.Milestones)...)

	return screen, nil
}

func buildHero(topic, imageKeyword string) Component {
	return Component{Type: TypeHero, ID: "hero", Props: map[string]any{"title": topic, "imageKeyword": imageKeyword}}
}

func buildProgressBar(percent int) Component {
	return Component{Type: TypeProgressBar, ID: "progress", Props: map[string]any{"percentage": percent}}
}

func buildCards(milestones []store.MilestoneRecord) []Component {
	var cards []Component
	for _, milestone := range milestones {
		cards = append(cards, buildCard(milestone.Title, milestone.Order))
	}
	return cards
}

func buildCard(title string, order int) Component {
	return Component{Type: TypeCard, ID: fmt.Sprintf("card-%d", order), Props: map[string]any{"title": title, "order": order}}
}

func buildSkeletonLoader() Component {
	return Component{Type: TypeSkeletonLoader, ID: "content-area"}
}

// BuildInitScreen builds the POST SSE "init" event tree.
func BuildInitScreen(topic, seedColor, heroKeyword string) (Component, error) {
	vars, err := PaletteFromSeed(seedColor)
	if err != nil {
		return Component{}, err
	}

	screen := Component{Type: TypeScreen, ID: "journey-screen", Props: map[string]any{"style": vars}}
	screen.Children = append(screen.Children, buildHero(topic, heroKeyword))
	screen.Children = append(screen.Children, buildSkeletonLoader())
	return screen, nil
}

// BuildContentUpdate builds the POST SSE "update" REPLACE payload.
func BuildContentUpdate(milestones []store.MilestoneRecord) UpdatePayload {
	return UpdatePayload{Action: "REPLACE", TargetID: "content-area", Components: buildCards(milestones)}
}
