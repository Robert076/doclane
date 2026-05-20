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

func (repo *AuditLogRepo) GetNotificationsForUser(ctx context.Context, userID int, departmentID *int, isAdmin bool, limit int) ([]events.Event, error) {
	var query string
	var args []any

	if isAdmin {
		query = `
			SELECT a.id, a.event_type, a.actor_id,
			       u.first_name AS actor_first_name, u.last_name AS actor_last_name,
			       a.resource_type, a.resource_id, a.metadata, a.occurred_at
			FROM audit_log a
			LEFT JOIN users u ON u.id = a.actor_id
			WHERE a.actor_id != $1 OR a.actor_id IS NULL
			ORDER BY a.occurred_at DESC
			LIMIT $2
		`
		args = []any{userID, limit}
	} else if departmentID != nil {
		// Professional: events on requests in their department, not triggered by them
		query = `
		SELECT a.id, a.event_type, a.actor_id,
			u.first_name AS actor_first_name, u.last_name AS actor_last_name,
			a.resource_type, a.resource_id, a.metadata, a.occurred_at
		FROM audit_log a
		LEFT JOIN users u ON u.id = a.actor_id
		INNER JOIN document_requests r ON 
			(a.resource_type = 'request' AND a.resource_id = r.id)
			OR (a.resource_type = 'document' AND (a.metadata->>'request_id')::int = r.id)
		WHERE r.department_id = $1
		AND (a.actor_id != $2 OR a.actor_id IS NULL)
		ORDER BY a.occurred_at DESC
		LIMIT $3
		`
		args = []any{*departmentID, userID, limit}
	} else {
		// Citizen: events on their own requests, not triggered by them
		query = `
		SELECT a.id, a.event_type, a.actor_id,
			u.first_name AS actor_first_name, u.last_name AS actor_last_name,
			a.resource_type, a.resource_id, a.metadata, a.occurred_at
		FROM audit_log a
		LEFT JOIN users u ON u.id = a.actor_id
		INNER JOIN document_requests r ON 
			(a.resource_type = 'request' AND a.resource_id = r.id)
			OR (a.resource_type = 'document' AND (a.metadata->>'request_id')::int = r.id)
		WHERE r.assignee = $1
		AND (a.actor_id != $1 OR a.actor_id IS NULL)
		ORDER BY a.occurred_at DESC
		LIMIT $2
		`
		args = []any{userID, limit}
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []events.Event
	for rows.Next() {
		var e events.Event
		var actorID *int
		var actorFirstName *string
		var actorLastName *string
		var metadata []byte

		if err := rows.Scan(
			&e.ID,
			&e.Type,
			&actorID,
			&actorFirstName,
			&actorLastName,
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
		e.ActorFirstName = actorFirstName
		e.ActorLastName = actorLastName

		if metadata != nil {
			if err := json.Unmarshal(metadata, &e.Metadata); err != nil {
				return nil, err
			}
		}

		result = append(result, e)
	}

	return result, rows.Err()
}

func (repo *AuditLogRepo) GetByResource(ctx context.Context, resourceType string, resourceID int) ([]events.Event, error) {
	query := `
		SELECT a.id, a.event_type, a.actor_id,
			u.first_name AS actor_first_name, u.last_name AS actor_last_name,
			a.resource_type, a.resource_id, a.metadata, a.occurred_at
		FROM audit_log a
		LEFT JOIN users u ON u.id = a.actor_id
		WHERE (a.resource_type = $1 AND a.resource_id = $2)
		OR (a.resource_type = 'document' AND a.metadata->>'request_id' = $3)
		ORDER BY a.occurred_at ASC
	`

	rows, err := repo.db.QueryContext(ctx, query, resourceType, resourceID, resourceID)
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
			&e.ID,
			&e.Type,
			&actorID,
			&e.ActorFirstName,
			&e.ActorLastName,
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
