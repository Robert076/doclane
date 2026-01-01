package repositories

import (
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) AddUser(user models.User) (int, error) {
	var id int

	err := repo.db.QueryRow(
		`INSERT INTO users (email, password_hash, role, professional_id, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		user.Email,
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

func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := repo.db.QueryRow(
		`SELECT id, email, password_hash, role, professional_id, is_active, created_at, updated_at
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.ProfessionalID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, err
		}

		return models.User{}, err
	}

	return user, nil
}
