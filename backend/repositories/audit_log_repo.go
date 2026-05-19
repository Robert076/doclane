package repositories

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Robert076/doclane/backend/events"
)

type AuditLogRepo struct {
	db *sql.DB
}

func NewAuditLogRepo(db *sql.DB) *AuditLogRepo {
	return &AuditLogRepo{db: db}
}

func (repo *AuditLogRepo) LogEvent(ctx context.Context, event events.Event) error {
	metadata, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO audit_log (event_type, actor_id, resource_type, resource_id, metadata, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	var actorID *int
	if event.ActorID != 0 {
		actorID = &event.ActorID
	}

	_, err = repo.db.ExecContext(ctx, query,
		event.Type,
		actorID,
		event.ResourceType,
		event.ResourceID,
		metadata,
		event.OccurredAt,
	)
	return err
}

func (repo *AuditLogRepo) GetByResource(ctx context.Context, resourceType string, resourceID int) ([]events.Event, error) {
	query := `
		SELECT event_type, actor_id, resource_type, resource_id, metadata, occurred_at
		FROM audit_log
		WHERE resource_type = $1 AND resource_id = $2
		ORDER BY occurred_at ASC
	`

	rows, err := repo.db.QueryContext(ctx, query, resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []events.Event
	for rows.Next() {
		var e events.Event
		var actorID *int
		var metadata []byte

		if err := rows.Scan(
			&e.Type,
			&actorID,
			&e.ResourceType,
			&e.ResourceID,
			&metadata,
			&e.OccurredAt,
		); err != nil {
			return nil, err
		}

		if actorID != nil {
			e.ActorID = *actorID
		}

		if metadata != nil {
			if err := json.Unmarshal(metadata, &e.Metadata); err != nil {
				return nil, err
			}
		}

		result = append(result, e)
	}

	return result, rows.Err()
}
