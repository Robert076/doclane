package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Robert076/doclane/backend/models"
)

type TagRepo struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) GetTags(ctx context.Context) ([]models.Tag, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, color, created_at FROM tags ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]models.Tag, 0)
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *TagRepo) GetTagByID(ctx context.Context, id int) (models.Tag, error) {
	var t models.Tag
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, color, created_at FROM tags WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt)
	return t, err
}

func (r *TagRepo) GetTagsByTemplateID(ctx context.Context, templateID int) ([]models.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.name, t.color, t.created_at
		FROM tags t
		JOIN template_tags tt ON tt.tag_id = t.id
		WHERE tt.template_id = $1
		ORDER BY t.name ASC
	`, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]models.Tag, 0)
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *TagRepo) CreateTag(ctx context.Context, dto models.TagDTOCreate) (models.Tag, error) {
	var t models.Tag
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tags (name, color) VALUES ($1, $2) RETURNING id, name, color, created_at`,
		dto.Name, dto.Color,
	).Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt)
	return t, err
}

func (r *TagRepo) UpdateTag(ctx context.Context, id int, dto models.TagDTOUpdate) (models.Tag, error) {
	var t models.Tag
	err := r.db.QueryRowContext(ctx, `
		UPDATE tags SET name = $1, color = $2
		WHERE id = $3
		RETURNING id, name, color, created_at
	`, dto.Name, dto.Color, id).Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt)
	return t, err
}

func (r *TagRepo) DeleteTag(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tags WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TagRepo) SetTemplateTags(ctx context.Context, templateID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		_, err := r.db.ExecContext(ctx,
			`DELETE FROM template_tags WHERE template_id = $1`, templateID,
		)
		return err
	}

	placeholders := make([]string, len(tagIDs))
	args := make([]interface{}, len(tagIDs)+1)
	args[0] = templateID
	for i, id := range tagIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}

	deleteQuery := fmt.Sprintf(
		`DELETE FROM template_tags WHERE template_id = $1 AND tag_id NOT IN (%s)`,
		strings.Join(placeholders, ", "),
	)
	if _, err := r.db.ExecContext(ctx, deleteQuery, args...); err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		if _, err := r.db.ExecContext(ctx,
			`INSERT INTO template_tags (template_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			templateID, tagID,
		); err != nil {
			return err
		}
	}
	return nil
}
