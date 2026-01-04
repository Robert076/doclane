package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DocumentService struct {
	documentRepo repositories.IDocumentRepository
	userRepo     repositories.IUserRepository
	s3Client     *s3.Client
	bucket       string
}

func NewDocumentService(documentRepo repositories.IDocumentRepository, userRepo repositories.IUserRepository, s3Client *s3.Client, bucket string) *DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
		userRepo:     userRepo,
		s3Client:     s3Client,
		bucket:       bucket,
	}
}

func (service *DocumentService) AddDocumentRequest(
	ctx context.Context,
	jwtUserId int,
	professionalId int,
	clientId int,
	title string,
	description *string,
	dueDate *time.Time,
) (int, error) {
	if jwtUserId != professionalId {
		return 0, errors.ErrForbidden{Msg: fmt.Sprintf("User with id %v is not allowed to add request to user with id %v", jwtUserId, clientId)}
	}

	client, err := service.userRepo.GetUserByID(ctx, clientId)
	if err != nil {
		return 0, errors.ErrNotFound{Msg: "Client not found."}
	}

	if client.ProfessionalID == nil || *client.ProfessionalID != strconv.Itoa(professionalId) {
		return 0, errors.ErrForbidden{Msg: "This client is not assigned to you."}
	}

	req := models.DocumentRequest{ProfessionalID: professionalId, ClientID: clientId, Title: title, Description: description, DueDate: dueDate}
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.Status = "pending"

	id, err := service.documentRepo.AddDocumentRequest(ctx, req)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (service *DocumentService) GetDocumentRequestByID(
	ctx context.Context,
	jwtUserId int,
	id int,
) (models.DocumentRequest, error) {
	req, err := service.documentRepo.GetDocumentRequestByID(ctx, id)
	if err != nil {
		return models.DocumentRequest{}, err
	}

	if req.ProfessionalID != jwtUserId && req.ClientID != jwtUserId {
		return models.DocumentRequest{}, errors.ErrForbidden{Msg: fmt.Sprintf("User with id %v is not allowed to access document request with id %v", jwtUserId, req.ID)}
	}

	return req, nil
}

func (service *DocumentService) GetDocumentRequestsByProfessional(
	ctx context.Context,
	jwtUserId int,
	professionalId int,
) ([]models.DocumentRequest, error) {
	if jwtUserId != professionalId {
		return nil, errors.ErrForbidden{Msg: fmt.Sprintf("User with id %v is not allowed to access document requests from professional with id %v", jwtUserId, professionalId)}
	}

	reqs, err := service.documentRepo.GetDocumentRequestsByProfessional(ctx, professionalId)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (service *DocumentService) GetDocumentRequestsByClient(
	ctx context.Context,
	jwtUserId int,
	clientId int,
) ([]models.DocumentRequest, error) {
	client, err := service.userRepo.GetUserByID(ctx, clientId)
	if err != nil {
		return nil, err
	}

	if !client.IsActive {
		return nil, errors.ErrForbidden{Msg: "This client account is deactivated."}
	}

	isOwner := clientId == jwtUserId
	isAssignedProfessional := false

	if client.ProfessionalID != nil {
		profId, _ := strconv.Atoi(*client.ProfessionalID)
		if profId == jwtUserId {
			isAssignedProfessional = true
		}
	}

	if !isOwner && !isAssignedProfessional {
		return nil, errors.ErrForbidden{Msg: "You do not have permission to view these requests."}
	}

	reqs, err := service.documentRepo.GetDocumentRequestsByClient(ctx, clientId)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (service *DocumentService) UpdateDocumentRequestStatus(
	ctx context.Context,
	jwtUserId int,
	id int,
	status string,
) error {
	req, err := service.documentRepo.GetDocumentRequestByID(ctx, id)
	if err != nil {
		return err
	}

	if req.ProfessionalID != jwtUserId && req.ClientID != jwtUserId {
		return errors.ErrForbidden{Msg: "You are not authorized to update this request status."}
	}

	validStatuses := map[string]bool{
		"pending":  true,
		"uploaded": true,
		"overdue":  true,
	}

	if !validStatuses[status] {
		return errors.ErrBadRequest{Msg: "Invalid status. Allowed: 'pending', 'uploaded', 'overdue'."}
	}

	if err := service.documentRepo.UpdateDocumentRequestStatus(ctx, id, status); err != nil {
		return err
	}

	return nil
}

func (service *DocumentService) AddDocumentFile(
	ctx context.Context,
	userId int,
	requestID int,
	fileName string,
	fileSize int64,
	contentType string,
	content io.Reader,
) (int, error) {
	docReq, err := service.documentRepo.GetDocumentRequestByID(ctx, requestID)
	if err != nil {
		return 0, errors.ErrNotFound{Msg: fmt.Sprintf("Document request not found. %v", err)}
	}

	if docReq.ClientID != userId && docReq.ProfessionalID != userId {
		return 0, errors.ErrForbidden{Msg: fmt.Sprintf("User with id %v is not allowed to modify document request with id %v.", userId, requestID)}
	}

	cleanFileName := filepath.Base(fileName)
	s3Key := fmt.Sprintf("requests/%d/%s", requestID, cleanFileName)

	result, err := service.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(service.bucket),
		Key:         aws.String(s3Key),
		Body:        content,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return 0, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to upload to S3. %v", err)}
	}

	fileModel := models.DocumentFile{
		DocumentRequestID: requestID,
		FileName:          cleanFileName,
		FilePath:          s3Key,
		MimeType:          &contentType,
		FileSize:          &fileSize,
		S3VersionID:       result.VersionId,
		UploadedAt:        time.Now(),
	}

	id, err := service.documentRepo.AddDocumentFile(ctx, fileModel)
	if err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, cleanupErr := service.s3Client.DeleteObject(cleanupCtx, &s3.DeleteObjectInput{
			Bucket:    aws.String(service.bucket),
			Key:       aws.String(s3Key),
			VersionId: result.VersionId,
		})

		if cleanupErr != nil {
			return 0, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to cleanup S3 object after DB failure. Key: %s, Version: %s, Error: %v\n", s3Key, *result.VersionId, cleanupErr)}
		}

		return 0, errors.ErrInternalServerError{Msg: fmt.Sprintf("Metadata save failed, file removed from storage. %v", err)}
	}

	return id, nil
}

func (service *DocumentService) GetFilesByRequest() {

}
