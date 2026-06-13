//go:build integration

// Package repositories integration tests run against a real PostgreSQL instance
// started in a throwaway container. They are gated behind the "integration"
// build tag so the default `go test ./...` stays fast and dependency-free.
//
// Run with:
//
//	go test -tags=integration ./repositories/...
//
// Requires a working Docker daemon (provided by GitHub-hosted runners).
package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// testDB is the shared connection pool for the lifetime of the test binary.
// It is initialised once in TestMain and reused by every test, which keeps the
// suite fast: one container start, not one per test.
var testDB *sql.DB

// schemaPath points at the canonical schema the application ships with, so the
// tests exercise the exact tables, constraints, and defaults used in production.
const schemaPath = "../../init.sql"

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, db, err := startPostgres(ctx)
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}
	testDB = db

	code := m.Run()

	// Best-effort cleanup; the container is ryuk-reaped even if this fails.
	_ = db.Close()
	if err := container.Terminate(ctx); err != nil {
		log.Printf("could not terminate postgres container: %v", err)
	}

	os.Exit(code)
}

// startPostgres spins up a Postgres container matching the production image,
// applies the application schema, and returns a ready-to-use connection pool.
func startPostgres(ctx context.Context) (*postgres.PostgresContainer, *sql.DB, error) {
	schema, err := filepath.Abs(schemaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("resolving schema path: %w", err)
	}

	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithInitScripts(schema),
		postgres.WithDatabase("doclane"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("running container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("building connection string: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("opening db: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("pinging db: %w", err)
	}

	return container, db, nil
}

// resetDB truncates every application table and restarts identity sequences so
// each test starts from a known-empty state. Register it with t.Cleanup or call
// it at the top of a test. CASCADE handles the foreign-key ordering for us.
func resetDB(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`
		TRUNCATE
			audit_log,
			template_tags,
			tags,
			document_comments,
			document_files,
			expected_documents,
			document_requests,
			expected_document_templates,
			document_request_templates,
			invitation_codes,
			users,
			departments
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to reset database: %v", err)
	}
}
