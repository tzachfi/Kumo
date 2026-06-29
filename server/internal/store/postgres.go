package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tzachfi/kumo/server/internal/domain"
)

// PostgresStore implements JourneyStore against PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

// NewPostgres returns a store backed by the given connection pool.
func NewPostgres(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// GetJourneyByID loads a journey and its milestones, mapping to JourneyRecord.
func (s *PostgresStore) GetJourneyByID(ctx context.Context, id string) (*JourneyRecord, error) {
	journeyID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrNotFound
	}

	var title string
	var progressPct int
	var config json.RawMessage

	row := s.pool.QueryRow(ctx,
		"SELECT title, progress_pct, config FROM journeys WHERE id = $1",
		journeyID,
	)
	if err := row.Scan(&title, &progressPct, &config); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("store: get journey: %w", err)
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id, title, sequence_order
		 FROM milestones
		 WHERE journey_id = $1
		 ORDER BY sequence_order`,
		journeyID,
	)
	if err != nil {
		return nil, fmt.Errorf("store: list milestones: %w", err)
	}
	defer rows.Close()

	var milestones []domain.Milestone
	for rows.Next() {
		var m domain.Milestone
		if err := rows.Scan(&m.ID, &m.Title, &m.Order); err != nil {
			return nil, fmt.Errorf("store: scan milestone: %w", err)
		}
		milestones = append(milestones, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: iterate milestones: %w", err)
	}

	rec := journeyToRecord(journeyID, title, progressPct, config, milestones)
	return &rec, nil
}

// CreateJourney inserts a journey row.
func (s *PostgresStore) CreateJourney(ctx context.Context, j *domain.Journey) error {
	if j == nil {
		return ErrBadInput
	}

	config := j.Config
	if len(config) == 0 {
		config = json.RawMessage("{}")
	}

	_, err := s.pool.Exec(ctx,
		`INSERT INTO journeys (id, user_id, title, state, deadline, config, progress_pct)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		j.ID, j.UserID, j.Title, j.State, j.Deadline, config, 0,
	)
	if err != nil {
		return fmt.Errorf("store: create journey: %w", err)
	}
	return nil
}

// SaveMilestones inserts milestones and updates journey state in one transaction.
func (s *PostgresStore) SaveMilestones(ctx context.Context, journeyID uuid.UUID, ms []domain.Milestone, state domain.JourneyState) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("store: begin save milestones: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, m := range ms {
		_, err = tx.Exec(ctx,
			`INSERT INTO milestones (id, journey_id, title, sequence_order, state)
			 VALUES ($1, $2, $3, $4, $5)`,
			m.ID, journeyID, m.Title, m.Order, m.State,
		)
		if err != nil {
			return fmt.Errorf("store: insert milestone: %w", err)
		}
	}

	_, err = tx.Exec(ctx,
		`UPDATE journeys SET state = $1, progress_pct = 0 WHERE id = $2`,
		state, journeyID,
	)
	if err != nil {
		return fmt.Errorf("store: update journey state: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("store: commit save milestones: %w", err)
	}
	return nil
}

var _ JourneyStore = (*PostgresStore)(nil)
