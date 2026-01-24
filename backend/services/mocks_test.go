package services

import (
	"context"

	"github.com/Robert076/doclane/backend/models"
)

type MockUserRepo struct {
	GetUserByIDFunc              func(ctx context.Context, id int) (models.User, error)
	GetUserByEmailFunc           func(ctx context.Context, email string) (models.User, error)
	GetUsersByProfessionalIDFunc func(ctx context.Context, professionalID int, limit *int, offset *int) ([]models.User, error)
	GetUsersFunc                 func(ctx context.Context, limit *int, offset *int, orderBy *string, order *string) ([]models.User, error)
	AddUserFunc                  func(ctx context.Context, user models.User) (int, error)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return models.User{}, nil
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return models.User{}, nil
}

func (m *MockUserRepo) GetUsersByProfessionalID(
	ctx context.Context,
	professionalID int,
	limit *int,
	offset *int,
) ([]models.User, error) {
	if m.GetUsersByProfessionalIDFunc != nil {
		return m.GetUsersByProfessionalIDFunc(ctx, professionalID, limit, offset)
	}
	return []models.User{}, nil
}

func (m *MockUserRepo) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
	orderBy *string,
	order *string,
) ([]models.User, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(ctx, limit, offset, orderBy, order)
	}
	return []models.User{}, nil
}

func (m *MockUserRepo) AddUser(ctx context.Context, user models.User) (int, error) {
	if m.AddUserFunc != nil {
		return m.AddUserFunc(ctx, user)
	}
	return 0, nil
}

type MockDocRepo struct {
	AddDocumentRequestFunc                func(ctx context.Context, req models.DocumentRequest) (int, error)
	GetDocumentRequestByIdFunc            func(ctx context.Context, id int) (models.DocumentRequestDTO, error)
	GetDocumentRequestsByProfessionalFunc func(ctx context.Context, id int) ([]models.DocumentRequestDTO, error)
	GetDocumentRequestsByClientFunc       func(ctx context.Context, id int) ([]models.DocumentRequestDTO, error)
	GetFilesByRequestFunc                 func(ctx context.Context, id int) ([]models.DocumentFile, error)
}

func (m *MockDocRepo) AddDocumentRequest(ctx context.Context, req models.DocumentRequest) (int, error) {
	if m.AddDocumentRequestFunc != nil {
		return m.AddDocumentRequestFunc(ctx, req)
	}
	return 0, nil
}

func (m *MockDocRepo) GetDocumentRequestByID(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
	if m.GetDocumentRequestByIdFunc != nil {
		return m.GetDocumentRequestByIdFunc(ctx, id)
	}
	return models.DocumentRequestDTO{}, nil
}
func (m *MockDocRepo) GetDocumentRequestsByProfessional(ctx context.Context, id int) ([]models.DocumentRequestDTO, error) {
	if m.GetDocumentRequestsByProfessionalFunc != nil {
		return m.GetDocumentRequestsByProfessionalFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockDocRepo) GetDocumentRequestsByClient(ctx context.Context, id int) ([]models.DocumentRequestDTO, error) {
	if m.GetDocumentRequestsByClientFunc != nil {
		return m.GetDocumentRequestsByClientFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockDocRepo) UpdateDocumentRequestStatus(ctx context.Context, id int, status string) error {
	return nil
}
func (m *MockDocRepo) AddDocumentFile(ctx context.Context, file models.DocumentFile) (int, error) {
	return 0, nil
}
func (m *MockDocRepo) GetFilesByRequest(ctx context.Context, id int) ([]models.DocumentFile, error) {
	if m.GetFilesByRequestFunc != nil {
		return m.GetFilesByRequestFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockDocRepo) GetFileByID(ctx context.Context, id int) (models.DocumentFile, error) {
	return models.DocumentFile{}, nil
}
