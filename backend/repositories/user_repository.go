package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
	orderBy *string,
	order *string,
	search *string,
) ([]models.User, error) {
	users := []models.User{}

	// coloane permise pentru ORDER BY (whitelist)
	allowedOrderBy := map[string]string{
		"id":         "id",
		"email":      "email",
		"created_at": "created_at",
		"updated_at": "updated_at",
		"role":       "role",
	}

	query := `
        SELECT id, email, first_name, last_name, password_hash, role, professional_id, is_active, created_at, updated_at
        FROM users
    `

	args := []interface{}{}
	argIndex := 1

	// WHERE clause for search
	if search != nil && *search != "" {
		searchPattern := "%" + strings.ToLower(*search) + "%"
		query += ` WHERE (
            LOWER(email) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(first_name) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(last_name) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(first_name || ' ' || last_name) LIKE $` + strconv.Itoa(argIndex) + `
        )`
		args = append(args, searchPattern)
		argIndex++
	}

	// ORDER BY
	if orderBy != nil {
		column, ok := allowedOrderBy[*orderBy]
		if ok {
			direction := "ASC"
			if order != nil && (*order == "asc" || *order == "desc") {
				direction = strings.ToUpper(*order)
			}
			query += " ORDER BY " + column + " " + direction
		}
	}

	// LIMIT
	if limit != nil {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		args = append(args, *limit)
		argIndex++
	}

	// OFFSET
	if offset != nil {
		query += " OFFSET $" + strconv.Itoa(argIndex)
		args = append(args, *offset)
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PasswordHash,
			&user.Role,
			&user.ProfessionalID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var user models.User

	err := repo.db.QueryRowContext(ctx,
		`SELECT id, email, first_name, last_name, password_hash, role, professional_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Role,
		&user.ProfessionalID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.ErrNotFound{Msg: "User not found."}
		}

		return models.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := repo.db.QueryRowContext(ctx,
		`SELECT id, email, first_name, last_name, password_hash, role, professional_id, is_active, created_at, updated_at
			FROM users
			WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Role,
		&user.ProfessionalID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.ErrNotFound{Msg: "User not found."}
		}

		return models.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) GetUsersByProfessionalID(
	ctx context.Context,
	professionalID int,
	limit *int,
	offset *int,
) ([]models.User, error) {
	users := []models.User{}

	query := `
        SELECT id, email, first_name, last_name, password_hash, role, professional_id, is_active, created_at, updated_at
        FROM users
        WHERE professional_id = $1
    `

	args := []interface{}{professionalID}
	argIndex := 2

	if limit != nil {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		args = append(args, *limit)
		argIndex++
	}

	if offset != nil {
		query += " OFFSET $" + strconv.Itoa(argIndex)
		args = append(args, *offset)
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PasswordHash,
			&user.Role,
			&user.ProfessionalID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) AddUser(ctx context.Context, user models.User) (int, error) {
	var id int

	err := repo.db.QueryRowContext(ctx,
		`INSERT INTO users (email, first_name, last_name, password_hash, role, professional_id, is_active, created_at, updated_at)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				 RETURNING id`,
		user.Email,
		&user.FirstName,
		&user.LastName,
		user.PasswordHash,
		user.Role,
		user.ProfessionalID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
