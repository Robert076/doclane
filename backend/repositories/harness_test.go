//go:build integration

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

var testDB *sql.DB

const schemaPath = "../../init.sql"

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, db, err := startPostgres(ctx)
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}
	testDB = db

	code := m.Run()

	_ = db.Close()
	if err := container.Terminate(ctx); err != nil {
		log.Printf("could not terminate postgres container: %v", err)
	}

	os.Exit(code)
}

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

// resetDB clears all tables so each test starts from an empty state.
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
