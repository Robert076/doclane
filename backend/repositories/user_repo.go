package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
	orderBy *string,
	order *string,
	search *string,
) ([]models.User, error) {
	allowedOrderBy := map[string]string{
		"id":         "id",
		"email":      "email",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	query := `
		SELECT id, email, first_name, last_name, password_hash, role, department_id, is_active, last_notified, created_at, updated_at, phone, street, locality
		FROM users
	`

	args := []interface{}{}
	argIndex := 1

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

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PasswordHash,
			&user.Role,
			&user.DepartmentID,
			&user.IsActive,
			&user.LastNotified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Phone,
			&user.Street,
			&user.Locality,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (repo *UserRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	query := `
		SELECT id, email, first_name, last_name, password_hash, role, department_id, is_active, last_notified, created_at, updated_at, phone, street, locality
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Role,
		&user.DepartmentID,
		&user.IsActive,
		&user.LastNotified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Phone,
		&user.Street,
		&user.Locality,
	)
	if err == sql.ErrNoRows {
		return models.User{}, errors.ErrNotFound{Msg: "User not found."}
	}
	return user, err
}

func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := `
		SELECT id, email, first_name, last_name, password_hash, role, department_id, is_active, last_notified, created_at, updated_at, phone, street, locality
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := repo.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Role,
		&user.DepartmentID,
		&user.IsActive,
		&user.LastNotified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Phone,
		&user.Street,
		&user.Locality,
	)
	if err == sql.ErrNoRows {
		return models.User{}, errors.ErrNotFound{Msg: "User not found."}
	}
	return user, err
}

func (repo *UserRepo) GetUsersByDepartment(ctx context.Context, departmentID int) ([]models.User, error) {
	query := `
		SELECT id, email, first_name, last_name, password_hash, role, department_id, is_active, last_notified, created_at, updated_at, phone, street, locality
		FROM users
		WHERE department_id = $1
	`
	rows, err := repo.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PasswordHash,
			&user.Role,
			&user.DepartmentID,
			&user.IsActive,
			&user.LastNotified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Phone,
			&user.Street,
			&user.Locality,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (repo *UserRepo) AddUser(ctx context.Context, user models.User) (int, error) {
	query := `
		INSERT INTO users (email, first_name, last_name, password_hash, role, department_id, is_active, last_notified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var id int
	err := repo.db.QueryRowContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PasswordHash,
		user.Role,
		user.DepartmentID,
		user.IsActive,
		user.LastNotified,
	).Scan(&id)
	return id, err
}

func (repo *UserRepo) UpdateUserProfile(ctx context.Context, userID int, dto models.UserProfilePatchDTO) error {
	query := `
		UPDATE users
		SET phone = $1, street = $2, locality = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := repo.db.ExecContext(ctx, query, dto.Phone, dto.Street, dto.Locality, userID)
	return err
}

func (repo *UserRepo) UpdateUserDepartment(ctx context.Context, userID int, departmentID int) error {
	query := `UPDATE users SET department_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := repo.db.ExecContext(ctx, query, departmentID, userID)
	return err
}

func (repo *UserRepo) DeactivateUser(ctx context.Context, userId int) error {
	query := `UPDATE users SET is_active=false WHERE id=$1`
	_, err := repo.db.ExecContext(ctx, query, userId)
	return err
}

func (repo *UserRepo) NotifyUser(ctx context.Context, userId int, time time.Time) error {
	query := `UPDATE users SET last_notified=$1 WHERE id=$2`
	_, err := repo.db.ExecContext(ctx, query, time, userId)
	return err
}

func (repo *UserRepo) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := repo.db.ExecContext(ctx, query, hashedPassword, userID)
	return err
}
