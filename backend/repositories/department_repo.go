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

func (r *DepartmentRepo) GetAllDepartments(ctx context.Context) ([]models.DepartmentDTORead, error) {
	query := `
        SELECT 
            d.id, 
            d.name, 
            d.created_at, 
            d.updated_at,
            COUNT(u.id) AS member_count
        FROM departments d
        LEFT JOIN users u ON u.department_id = d.id
        GROUP BY d.id, d.name, d.created_at, d.updated_at
        ORDER BY d.created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departments := make([]models.DepartmentDTORead, 0)
	for rows.Next() {
		var d models.DepartmentDTORead
		if err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt, &d.MemberCount); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, rows.Err()
}

func (r *DepartmentRepo) GetDepartmentByID(ctx context.Context, id int) (models.DepartmentDTORead, error) {
	var d models.DepartmentDTORead
	query := `
        SELECT 
            d.id, 
            d.name, 
            d.created_at, 
            d.updated_at,
            COUNT(u.id) AS member_count
        FROM departments d
        LEFT JOIN users u ON u.department_id = d.id
        WHERE d.id = $1
        GROUP BY d.id, d.name, d.created_at, d.updated_at
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.Name,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.MemberCount,
	)
	return d, err
}

func (r *DepartmentRepo) CreateDepartment(ctx context.Context, name string) (int, error) {
	var id int
	query := `INSERT INTO departments (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
	return id, err
}
