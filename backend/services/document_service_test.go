package services

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Robert076/doclane/backend/models"
)

func TestDocumentService_AddDocumentRequest(t *testing.T) {
	PROFESSIONAL_ID := 10
	CLIENT_ID := 5
	CLIENT_ID_STR := "5"
	NEW_REQUEST_ID := 100
	OTHER_PROFESSIONAL_ID := 99
	NON_EXISTENT_CLIENT_ID := 999
	VALID_TITLE := "Valid Title"
	INVALID_TITLE := "No"
	CLIENT_ROLE := "CLIENT"

	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)

	type args struct {
		jwtUserId   int
		clientId    int
		title       string
		description *string
		dueDate     *time.Time
	}

	tests := []struct {
		name          string
		args          args
		setupMocks    func(mUser *MockUserRepo, mDoc *MockDocRepo)
		expectedID    int
		expectError   bool
		errorContains string
	}{
		{
			name: "Success - Request created",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
				clientId:  CLIENT_ID,
				title:     VALID_TITLE,
				dueDate:   &futureDate,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					profiID := strconv.Itoa(PROFESSIONAL_ID)
					return models.User{
						ID:             CLIENT_ID_STR,
						Role:           CLIENT_ROLE,
						ProfessionalID: &profiID,
					}, nil
				}

				mDoc.AddDocumentRequestFunc = func(ctx context.Context, req models.DocumentRequest) (int, error) {
					return NEW_REQUEST_ID, nil
				}
			},
			expectedID:  NEW_REQUEST_ID,
			expectError: false,
		},

		{
			name: "Failure - Title too short",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
				clientId:  CLIENT_ID,
				title:     INVALID_TITLE,
				dueDate:   &futureDate,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
			},
			expectedID:    0,
			expectError:   true,
			errorContains: "Title must be between",
		},

		{
			name: "Failure - Due date in past",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
				clientId:  CLIENT_ID,
				title:     VALID_TITLE,
				dueDate:   &pastDate,
			},
			setupMocks:    func(mUser *MockUserRepo, mDoc *MockDocRepo) {},
			expectedID:    0,
			expectError:   true,
			errorContains: "past",
		},

		{
			name: "Failure - Client not found",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
				clientId:  NON_EXISTENT_CLIENT_ID,
				title:     VALID_TITLE,
				dueDate:   &futureDate,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					return models.User{}, errors.New("user not found")
				}
			},
			expectedID:    0,
			expectError:   true,
			errorContains: "Client not found",
		},

		{
			name: "Failure - Forbidden (Client belongs to another pro)",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
				clientId:  CLIENT_ID,
				title:     VALID_TITLE,
				dueDate:   &futureDate,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					otherProfiID := strconv.Itoa(OTHER_PROFESSIONAL_ID)
					return models.User{
						ID:             CLIENT_ID_STR,
						ProfessionalID: &otherProfiID,
					}, nil
				}
			},
			expectedID:    0,
			expectError:   true,
			errorContains: "not assigned to you",
		},

		{
			name: "Failure - Database Error on Save",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
				clientId:  CLIENT_ID,
				title:     VALID_TITLE,
				dueDate:   &futureDate,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					profiID := strconv.Itoa(PROFESSIONAL_ID)
					return models.User{ID: CLIENT_ID_STR, ProfessionalID: &profiID}, nil
				}

				mDoc.AddDocumentRequestFunc = func(ctx context.Context, req models.DocumentRequest) (int, error) {
					return 0, errors.New("connection timeout")
				}
			},
			expectedID:    0,
			expectError:   true,
			errorContains: "connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := &MockUserRepo{}
			mockDoc := &MockDocRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockUser, mockDoc)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			service := NewDocumentService(mockDoc, mockUser, nil, "", logger)

			gotID, err := service.AddDocumentRequest(
				context.Background(),
				tt.args.jwtUserId,
				tt.args.clientId,
				tt.args.title,
				tt.args.description,
				tt.args.dueDate,
			)

			if (err != nil) != tt.expectError {
				t.Errorf("AddDocumentRequest() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error message = %v, want to contain %v", err.Error(), tt.errorContains)
				}
			}

			if gotID != tt.expectedID {
				t.Errorf("AddDocumentRequest() = %v, want %v", gotID, tt.expectedID)
			}
		})
	}
}

func TestDocumentService_GetDocumentRequestById(t *testing.T) {
	DIFFERENT_CLIENT_ID := 999
	DIFFERENT_PROFESSIONAL_ID := 998
	CLIENT_ID := 10
	REQUEST_ID := 5
	NULL_REQUEST_ID := 0

	type args struct {
		jwtUserId int
		requestId int
	}

	tests := []struct {
		name          string
		args          args
		setupMocks    func(mUser *MockUserRepo, mDoc *MockDocRepo)
		expectedID    int
		expectError   bool
		errorContains string
	}{
		{
			name: "Success - the document belongs to the authenticated client",
			args: args{
				jwtUserId: CLIENT_ID,
				requestId: REQUEST_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:       REQUEST_ID,
						ClientID: CLIENT_ID,
					}, nil
				}
			},
			expectedID:    REQUEST_ID,
			expectError:   false,
			errorContains: "",
		},
		{
			name: "Success - the document belongs to the authenticated professional",
			args: args{
				jwtUserId: CLIENT_ID,
				requestId: REQUEST_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: CLIENT_ID,
					}, nil
				}
			},
			expectedID:    REQUEST_ID,
			expectError:   false,
			errorContains: "",
		},
		{
			name: "Failure - the document does not belong to the user",
			args: args{
				jwtUserId: CLIENT_ID,
				requestId: REQUEST_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: DIFFERENT_PROFESSIONAL_ID,
						ClientID:       DIFFERENT_CLIENT_ID,
					}, nil
				}
			},
			expectedID:    NULL_REQUEST_ID,
			expectError:   true,
			errorContains: "not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDoc := &MockDocRepo{}
			mockUser := &MockUserRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockUser, mockDoc)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewDocumentService(mockDoc, mockUser, nil, "", logger)

			got, err := service.GetDocumentRequestByID(context.Background(), tt.args.jwtUserId, tt.args.requestId)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error message mismatch. Want substring: %s, Got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if got.ID != tt.expectedID {
					t.Errorf("Expected ID %d, got %d", tt.expectedID, got.ID)
				}
			}
		})
	}
}

func TestDocumentService_GetDocumentRequestByProfessional(t *testing.T) {
	PROFESSIONAL_ID := 1
	PROFESSIONAL_ID_STR := "1"
	ROLE_PROFESSIONAL := "PROFESSIONAL"
	ROLE_CLIENT := "CLIENT"
	NON_PROFESSIONAL_ID_STR := "2"
	DOCUMENT_REQUEST_ID := 100

	type args struct {
		jwtUserId int
	}

	tests := []struct {
		name        string
		args        args
		setupMocks  func(mUser *MockUserRepo, mDoc *MockDocRepo)
		expectError bool
	}{
		{
			name: "Success - User is a professional and can therefore read his clients",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					return models.User{
						ID:   PROFESSIONAL_ID_STR,
						Role: ROLE_PROFESSIONAL,
					}, nil
				}
				mDoc.GetDocumentRequestsByProfessionalFunc = func(ctx context.Context, id int) ([]models.DocumentRequestDTO, error) {
					return []models.DocumentRequestDTO{
						{
							ID: DOCUMENT_REQUEST_ID,
						},
					}, nil
				}
			},
			expectError: false,
		},
		{
			name: "Failure - User is not a professional, should return forbidden",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					return models.User{
						ID:   NON_PROFESSIONAL_ID_STR,
						Role: ROLE_CLIENT,
					}, nil
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDoc := &MockDocRepo{}
			mockUser := &MockUserRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockUser, mockDoc)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewDocumentService(mockDoc, mockUser, nil, "", logger)

			_, err := service.GetDocumentRequestsByProfessional(context.Background(), tt.args.jwtUserId)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDocumentService_GetDocumentRequestsByClient(t *testing.T) {
	CLIENT_ID := 5
	CLIENT_ID_STR := "5"
	CLIENT_ROLE := "CLIENT"
	PROFESSIONAL_ID := 10
	PROFESSIONAL_ID_STR := "10"
	PROFESSIONAL_ROLE := "PROFESSIONAL"
	REQUEST_ID := 100

	type args struct {
		jwtUserId int
	}

	tests := []struct {
		name        string
		args        args
		setupMocks  func(mUser *MockUserRepo, mDoc *MockDocRepo)
		expectError bool
	}{
		{
			name: "Success - User is a client",
			args: args{
				jwtUserId: CLIENT_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					return models.User{
						ID:   CLIENT_ID_STR,
						Role: CLIENT_ROLE,
					}, nil
				}
				mDoc.GetDocumentRequestsByClientFunc = func(ctx context.Context, id int) ([]models.DocumentRequestDTO, error) {
					return []models.DocumentRequestDTO{
						{ID: REQUEST_ID},
					}, nil
				}
			},
			expectError: false,
		},
		{
			name: "Failure - User is not a client",
			args: args{
				jwtUserId: PROFESSIONAL_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mUser.GetUserByIDFunc = func(ctx context.Context, id int) (models.User, error) {
					return models.User{
						ID:   PROFESSIONAL_ID_STR,
						Role: PROFESSIONAL_ROLE,
					}, nil
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDoc := &MockDocRepo{}
			mockUser := &MockUserRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockUser, mockDoc)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewDocumentService(mockDoc, mockUser, nil, "", logger)

			_, err := service.GetDocumentRequestsByClient(context.Background(), tt.args.jwtUserId)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDocumentService_UpdateDocumentRequestStatus(t *testing.T) {
	USER_ID := 10
	REQUEST_ID := 100
	VALID_STATUS := "uploaded"
	INVALID_STATUS := "finished"
	UNAUTHORIZED_USER_ID := 99

	type args struct {
		jwtUserId int
		id        int
		status    string
	}

	tests := []struct {
		name        string
		args        args
		setupMocks  func(mUser *MockUserRepo, mDoc *MockDocRepo)
		expectError bool
	}{
		{
			name: "Success - Status updated by owner",
			args: args{
				jwtUserId: USER_ID,
				id:        REQUEST_ID,
				status:    VALID_STATUS,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: USER_ID,
					}, nil
				}
			},
			expectError: false,
		},
		{
			name: "Failure - Invalid status",
			args: args{
				jwtUserId: USER_ID,
				id:        REQUEST_ID,
				status:    INVALID_STATUS,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: USER_ID,
					}, nil
				}
			},
			expectError: true,
		},
		{
			name: "Failure - Unauthorized user",
			args: args{
				jwtUserId: UNAUTHORIZED_USER_ID,
				id:        REQUEST_ID,
				status:    VALID_STATUS,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: USER_ID,
					}, nil
				}
			},
			expectError: true,
		},
		{
			name: "Failure - Request not found",
			args: args{
				jwtUserId: USER_ID,
				id:        REQUEST_ID,
				status:    VALID_STATUS,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{}, errors.New("not found")
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDoc := &MockDocRepo{}
			mockUser := &MockUserRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockUser, mockDoc)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewDocumentService(mockDoc, mockUser, nil, "", logger)

			err := service.UpdateDocumentRequestStatus(context.Background(), tt.args.jwtUserId, tt.args.id, tt.args.status)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDocumentService_GetFilesByRequest(t *testing.T) {
	USER_ID := 10
	REQUEST_ID := 50
	UNAUTHORIZED_USER_ID := 99
	FILE_ID := 1

	type args struct {
		jwtUserId int
		requestId int
	}

	tests := []struct {
		name        string
		args        args
		setupMocks  func(mUser *MockUserRepo, mDoc *MockDocRepo)
		expectError bool
	}{
		{
			name: "Success - Retrieve files for authorized user",
			args: args{
				jwtUserId: USER_ID,
				requestId: REQUEST_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: USER_ID,
					}, nil
				}
				mDoc.GetFilesByRequestFunc = func(ctx context.Context, id int) ([]models.DocumentFile, error) {
					return []models.DocumentFile{
						{ID: FILE_ID},
					}, nil
				}
			},
			expectError: false,
		},
		{
			name: "Failure - Unauthorized access",
			args: args{
				jwtUserId: UNAUTHORIZED_USER_ID,
				requestId: REQUEST_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: USER_ID,
					}, nil
				}
			},
			expectError: true,
		},
		{
			name: "Failure - Request not found",
			args: args{
				jwtUserId: USER_ID,
				requestId: REQUEST_ID,
			},
			setupMocks: func(mUser *MockUserRepo, mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{}, errors.New("not found")
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDoc := &MockDocRepo{}
			mockUser := &MockUserRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockUser, mockDoc)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewDocumentService(mockDoc, mockUser, nil, "", logger)

			_, err := service.GetFilesByRequest(context.Background(), tt.args.jwtUserId, tt.args.requestId)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDocumentService_AddDocumentFile_Validation(t *testing.T) {
	USER_ID := 10
	REQUEST_ID := 50
	INVALID_EXT_FILE := "test.exe"
	TOO_LARGE_FILE := int64(25 * 1024 * 1024)
	VALID_SIZE := int64(1024)
	VALID_NAME := "test.pdf"
	VALID_TYPE := "application/pdf"
	UNAUTHORIZED_ID := 99

	type args struct {
		jwtUserId   int
		requestID   int
		fileName    string
		fileSize    int64
		contentType string
		content     io.Reader
	}

	tests := []struct {
		name        string
		args        args
		setupMocks  func(mDoc *MockDocRepo)
		expectError bool
	}{
		{
			name: "Failure - File extension not allowed",
			args: args{
				jwtUserId:   USER_ID,
				requestID:   REQUEST_ID,
				fileName:    INVALID_EXT_FILE,
				fileSize:    VALID_SIZE,
				contentType: VALID_TYPE,
				content:     strings.NewReader("content"),
			},
			setupMocks:  func(mDoc *MockDocRepo) {},
			expectError: true,
		},
		{
			name: "Failure - File too large",
			args: args{
				jwtUserId:   USER_ID,
				requestID:   REQUEST_ID,
				fileName:    VALID_NAME,
				fileSize:    TOO_LARGE_FILE,
				contentType: VALID_TYPE,
				content:     strings.NewReader("content"),
			},
			setupMocks:  func(mDoc *MockDocRepo) {},
			expectError: true,
		},
		{
			name: "Failure - Request not found",
			args: args{
				jwtUserId:   USER_ID,
				requestID:   REQUEST_ID,
				fileName:    VALID_NAME,
				fileSize:    VALID_SIZE,
				contentType: VALID_TYPE,
				content:     strings.NewReader("content"),
			},
			setupMocks: func(mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{}, errors.New("not found")
				}
			},
			expectError: true,
		},
		{
			name: "Failure - Unauthorized upload",
			args: args{
				jwtUserId:   UNAUTHORIZED_ID,
				requestID:   REQUEST_ID,
				fileName:    VALID_NAME,
				fileSize:    VALID_SIZE,
				contentType: VALID_TYPE,
				content:     strings.NewReader("content"),
			},
			setupMocks: func(mDoc *MockDocRepo) {
				mDoc.GetDocumentRequestByIdFunc = func(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
					return models.DocumentRequestDTO{
						ID:             REQUEST_ID,
						ProfessionalID: USER_ID,
					}, nil
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDoc := &MockDocRepo{}
			mockUser := &MockUserRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockDoc)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := NewDocumentService(mockDoc, mockUser, nil, "", logger)

			_, err := service.AddDocumentFile(
				context.Background(),
				tt.args.jwtUserId,
				tt.args.requestID,
				tt.args.fileName,
				tt.args.fileSize,
				tt.args.contentType,
				tt.args.content,
			)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
