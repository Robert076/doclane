package repositories

import "github.com/Robert076/doclane/backend/models"

type IUserRepository interface {
	GetUsers(limit *int, offset *int, orderBy *string, order *string) ([]models.User, error)
	AddUser(user models.User) (int, error)
	GetUserByEmail(email string) (models.User, error)
}

type IDocumentRepository interface {
	CreateDocumentRequest(req models.DocumentRequest) (int, error)
	GetDocumentRequestByID(id int) (models.DocumentRequest, error)
	GetDocumentRequestsByProfessional(professionalID int) ([]models.DocumentRequest, error)
	GetDocumentRequestsByClient(clientID int) ([]models.DocumentRequest, error)
	UpdateDocumentRequestStatus(id int, status string) error

	AddDocumentFile(file models.DocumentFile) (int, error)
	GetFilesByRequest(requestID int) ([]models.DocumentFile, error)
}
