package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/tzachfi/kumo/server/internal/api"
	"github.com/tzachfi/kumo/server/internal/domain"
	"github.com/tzachfi/kumo/server/internal/generate"
	"github.com/tzachfi/kumo/server/internal/journey"
	"github.com/tzachfi/kumo/server/internal/sdui"
	"github.com/tzachfi/kumo/server/internal/store"
)

// JourneyHandler serves journey SDUI endpoints.
type JourneyHandler struct {
	Store         store.JourneyStore
	DefaultUserID uuid.UUID
}

type generateRequest struct {
	Topic string `json:"topic"`
}

type themeResult struct {
	seedColor   string
	heroKeyword string
}

type journeyConfig struct {
	SeedColor   string `json:"seed_color"`
	HeroKeyword string `json:"hero_keyword"`
}

type doneEvent struct {
	ID string `json:"id"`
}

// Get handles GET /api/journey/{id}.
func (h *JourneyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rec, err := h.Store.GetJourneyByID(r.Context(), id)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "Journey not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	screen, err := sdui.BuildDashboardScreen(*rec)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(screen)
}

// Generate handles POST /api/journey/generate (SSE).
func (h *JourneyHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if req.Topic == "" {
		http.Error(w, "Topic is required", http.StatusBadRequest)
		return
	}

	setSSEHeaders(w)
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	topic := req.Topic
	var journeyID uuid.UUID

	themeCh := make(chan themeResult, 1)
	contentCh := make(chan []store.MilestoneRecord, 1)

	go func() {
		time.Sleep(1 * time.Second)
		seedColor, heroKeyword := generate.MockTheme(topic)
		themeCh <- themeResult{seedColor: seedColor, heroKeyword: heroKeyword}
	}()

	go func() {
		time.Sleep(4 * time.Second)
		contentCh <- generate.MockMilestones(topic)
	}()

	pending := 2

	for pending > 0 {
		select {
		case <-r.Context().Done():
			return
		case theme := <-themeCh:
			config, err := json.Marshal(journeyConfig{
				SeedColor:   theme.seedColor,
				HeroKeyword: theme.heroKeyword,
			})

			if err != nil {
				return
			}

			j := &domain.Journey{
				ID:       uuid.New(),
				Title:    topic,
				State:    domain.Initializing,
				UserID:   h.DefaultUserID,
				Deadline: time.Now().UTC().Add(90 * 24 * time.Hour),
				Config:   config,
			}

			err = h.Store.CreateJourney(r.Context(), j)

			if err != nil {
				return
			}

			journeyID = j.ID

			screen, err := sdui.BuildInitScreen(topic, theme.seedColor, theme.heroKeyword)
			if err != nil {
				return
			}
			if err := api.WriteSSE(w, flusher, "init", screen); err != nil {
				return
			}
			pending--
		case milestones := <-contentCh:
			var dm = make([]domain.Milestone, 0, len(milestones))

			for i, rec := range milestones {
				state := domain.MilestoneLocked
				if i == 0 {
					state = domain.MilestoneActive
				}
				dm = append(dm, domain.Milestone{
					Title: rec.Title,
					Order: rec.Order,
					State: state,
				})
			}

			j := &domain.Journey{
				ID:         journeyID,
				Milestones: dm,
			}

			req := domain.JourneyContext{
				UserID:   h.DefaultUserID,
				Ambition: topic,
				Deadline: time.Now().UTC().Add(90 * 24 * time.Hour),
			}

			journey.Assemble(j, req)

			err = journey.Validate(j)

			if err != nil {
				return
			}

			err = h.Store.SaveMilestones(r.Context(), journeyID, j.Milestones, domain.Active)

			if err != nil {
				return
			}

			screenUpdate := sdui.BuildContentUpdate(milestones)
			if err := api.WriteSSE(w, flusher, "update", screenUpdate); err != nil {
				return
			}
			pending--
		}
	}

	if err := api.WriteSSE(w, flusher, "done", doneEvent{ID: journeyID.String()}); err != nil {
		return
	}
}

func setSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}
