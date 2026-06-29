//go:build integration

package store_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/tzachfi/kumo/server/internal/domain"
	"github.com/tzachfi/kumo/server/internal/journey"
	"github.com/tzachfi/kumo/server/internal/store"
)

const devUserID = "00000000-0000-4000-8000-000000000001"

func TestPostgresStore_RoundTrip(t *testing.T) {
	st, ctx := newIntegrationStore(t)

	userID := uuid.MustParse(devUserID)
	config, err := json.Marshal(map[string]string{
		"seed_color":   "#ff5500",
		"hero_keyword": "run",
	})
	if err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().UTC().Add(90 * 24 * time.Hour)
	j := &domain.Journey{
		ID:       uuid.New(),
		Title:    "10km Training",
		State:    domain.Initializing,
		UserID:   userID,
		Deadline: deadline,
		Config:   config,
	}
	if err := st.CreateJourney(ctx, j); err != nil {
		t.Fatalf("CreateJourney: %v", err)
	}

	j.Milestones = []domain.Milestone{
		{Title: "Base Building", Order: 0, State: domain.MilestoneActive},
		{Title: "Race Prep", Order: 1, State: domain.MilestoneLocked},
	}
	journey.Assemble(j, domain.JourneyContext{
		UserID:   userID,
		Ambition: j.Title,
		Deadline: deadline,
	})
	if err := journey.Validate(j); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if err := st.SaveMilestones(ctx, j.ID, j.Milestones, domain.Active); err != nil {
		t.Fatalf("SaveMilestones: %v", err)
	}

	rec, err := st.GetJourneyByID(ctx, j.ID.String())
	if err != nil {
		t.Fatalf("GetJourneyByID: %v", err)
	}
	if rec.Topic != "10km Training" {
		t.Errorf("Topic = %q, want %q", rec.Topic, "10km Training")
	}
	if rec.SeedColor != "#ff5500" {
		t.Errorf("SeedColor = %q, want %q", rec.SeedColor, "#ff5500")
	}
	if len(rec.Milestones) != 2 {
		t.Fatalf("milestones len = %d, want 2", len(rec.Milestones))
	}
	if rec.Milestones[0].Order != 0 || rec.Milestones[1].Order != 1 {
		t.Errorf("milestone order = [%d, %d], want [0, 1]", rec.Milestones[0].Order, rec.Milestones[1].Order)
	}
}

func TestPostgresStore_GetJourneyByID_NotFound(t *testing.T) {
	st, ctx := newIntegrationStore(t)

	_, err := st.GetJourneyByID(ctx, "not-a-uuid")
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("bad id: err = %v, want ErrNotFound", err)
	}

	_, err = st.GetJourneyByID(ctx, uuid.New().String())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("unknown id: err = %v, want ErrNotFound", err)
	}
}

func newIntegrationStore(t *testing.T) (*store.PostgresStore, context.Context) {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("kumo"),
		postgres.WithUsername("kumo"),
		postgres.WithPassword("kumo"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("postgres.Run: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("container.Terminate: %v", err)
		}
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("ConnectionString: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := waitForPool(ctx, pool); err != nil {
		t.Fatalf("waitForPool: %v", err)
	}

	if err := runMigrations(t, connStr); err != nil {
		t.Fatalf("runMigrations: %v", err)
	}

	return store.NewPostgres(pool), ctx
}

func waitForPool(ctx context.Context, pool *pgxpool.Pool) error {
	deadline := time.Now().Add(30 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := pool.Ping(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(200 * time.Millisecond)
	}
	if lastErr != nil {
		return lastErr
	}
	return pool.Ping(ctx)
}

func runMigrations(t *testing.T, connStr string) error {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	migrationsDir := filepath.Join(filepath.Dir(file), "..", "..", "migrations")

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, migrationsDir)
}
