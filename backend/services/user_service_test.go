package services

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/Robert076/doclane/backend/models"
	appErrors "github.com/Robert076/doclane/backend/types/errors"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_AddUser(t *testing.T) {
	const (
		VALID_EMAIL = "test@test.com"
		VALID_PASS  = "Password123!"
		ROLE_CLIENT = "CLIENT"
		USER_ID     = 10
	)

	type args struct {
		email string
		pass  string
		role  string
	}

	tests := []struct {
		name          string
		args          args
		setupMocks    func(repo *MockUserRepo)
		expectError   bool
		errorContains string
	}{
		{
			name: "Success - user created",
			args: args{
				email: VALID_EMAIL,
				pass:  VALID_PASS,
				role:  ROLE_CLIENT,
			},
			setupMocks: func(repo *MockUserRepo) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{}, appErrors.ErrNotFound{}
				}

				repo.AddUserFunc = func(ctx context.Context, user models.User) (int, error) {
					return USER_ID, nil
				}
			},
			expectError: false,
		},

		{
			name: "Failure - invalid email",
			args: args{
				email: "invalid",
				pass:  VALID_PASS,
				role:  ROLE_CLIENT,
			},
			setupMocks:    func(repo *MockUserRepo) {},
			expectError:   true,
			errorContains: "Invalid email",
		},

		{
			name: "Failure - invalid role",
			args: args{
				email: VALID_EMAIL,
				pass:  VALID_PASS,
				role:  "ADMIN",
			},
			setupMocks:    func(repo *MockUserRepo) {},
			expectError:   true,
			errorContains: "Invalid role",
		},

		{
			name: "Failure - user already exists",
			args: args{
				email: VALID_EMAIL,
				pass:  VALID_PASS,
				role:  ROLE_CLIENT,
			},
			setupMocks: func(repo *MockUserRepo) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{Email: email}, nil
				}
			},
			expectError:   true,
			errorContains: "already exists",
		},

		{
			name: "Failure - DB error on AddUser",
			args: args{
				email: VALID_EMAIL,
				pass:  VALID_PASS,
				role:  ROLE_CLIENT,
			},
			setupMocks: func(repo *MockUserRepo) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{}, appErrors.ErrNotFound{}
				}

				repo.AddUserFunc = func(ctx context.Context, user models.User) (int, error) {
					return 0, errors.New("db down")
				}
			},
			expectError:   true,
			errorContains: "db down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepo{}
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewUserService(mockRepo, logger)

			_, err := service.AddUser(context.Background(), CreateUserParams{
				Email:    tt.args.email,
				Password: tt.args.pass,
				Role:     tt.args.role,
			})

			if (err != nil) != tt.expectError {
				t.Fatalf("expected error=%v, got=%v", tt.expectError, err)
			}

			if tt.expectError && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error = %v, want contains %v", err.Error(), tt.errorContains)
				}
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	const EMAIL = "login@test.com"
	const PASSWORD = "Secret123!"

	hash, _ := bcrypt.GenerateFromPassword([]byte(PASSWORD), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(repo *MockUserRepo)
		expectError   bool
		errorContains string
	}{
		{
			name:     "Success - login ok",
			email:    EMAIL,
			password: PASSWORD,
			setupMocks: func(repo *MockUserRepo) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{
						ID:           "1",
						Email:        EMAIL,
						PasswordHash: string(hash),
					}, nil
				}
			},
		},

		{
			name:     "Failure - email not found",
			email:    EMAIL,
			password: PASSWORD,
			setupMocks: func(repo *MockUserRepo) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{}, appErrors.ErrNotFound{}
				}
			},
			expectError:   true,
			errorContains: "Invalid email or password",
		},

		{
			name:     "Failure - wrong password",
			email:    EMAIL,
			password: "wrong",
			setupMocks: func(repo *MockUserRepo) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{
						ID:           "1",
						Email:        EMAIL,
						PasswordHash: string(hash),
					}, nil
				}
			},
			expectError:   true,
			errorContains: "Invalid email or password",
		},

		{
			name:     "Failure - DB error",
			email:    EMAIL,
			password: PASSWORD,
			setupMocks: func(repo *MockUserRepo) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{}, errors.New("db offline")
				}
			},
			expectError:   true,
			errorContains: "db offline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepo{}
			tt.setupMocks(mockRepo)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewUserService(mockRepo, logger)

			_, err := service.Login(context.Background(), LoginParams{
				Email:    tt.email,
				Password: tt.password,
			})

			if (err != nil) != tt.expectError {
				t.Fatalf("expected error=%v, got=%v", tt.expectError, err)
			}

			if tt.expectError && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error = %v, want contains %v", err.Error(), tt.errorContains)
				}
			}
		})
	}
}
