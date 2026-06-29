// Command server is the HTTP entrypoint for the Kumo SDUI API.
package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tzachfi/kumo/server/internal/api/handler"
	"github.com/tzachfi/kumo/server/internal/store"

	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envStore         = "KUMO_STORE"
	envDatabaseURL   = "DATABASE_URL"
	envDefaultUserID = "KUMO_DEFAULT_USER_ID"
)

func main() {

	mux := http.NewServeMux()

	var st store.JourneyStore

	envs := os.Getenv(envStore)

	switch envs {
	case "postgres":
		databaseURL := os.Getenv(envDatabaseURL)
		if databaseURL == "" {
			log.Fatal("DATABASE_URL is required when KUMO_STORE=postgres")
		}
		pool, err := pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			log.Fatal(err)
		}
		defer pool.Close()
		st = store.NewPostgres(pool)
		log.Println("using store: postgres")
	default:
		st = store.NewMock()
		log.Println("using store: mock")
	}
	defaultUserID, err := uuid.Parse(os.Getenv(envDefaultUserID))
	if err != nil {
		log.Fatal(err)
	}

	jh := &handler.JourneyHandler{Store: st, DefaultUserID: defaultUserID}

	mux.HandleFunc("GET /api/journey/{id}", jh.Get)
	mux.HandleFunc("POST /api/journey/generate", jh.Generate)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down....")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("shut down")

}
