package repositories

import (
	"context"
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type DepartmentRepo struct {
	db *sql.DB
}

func NewDepartmentRepo(db *sql.DB) *DepartmentRepo {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM departments
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departments := make([]models.Department, 0)
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, rows.Err()
}

func (r *DepartmentRepo) GetDepartmentByID(ctx context.Context, id int) (models.Department, error) {
	var d models.Department
	query := `SELECT id, name, created_at, updated_at FROM departments WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt)
	return d, err
}

func (r *DepartmentRepo) CreateDepartment(ctx context.Context, name string) (int, error) {
	var id int
	query := `INSERT INTO departments (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
	return id, err
}
