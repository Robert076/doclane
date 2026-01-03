package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DocumentService struct {
	repo     repositories.IDocumentRepository
	s3Client *s3.Client
	bucket   string
}

func NewDocumentService(repo repositories.IDocumentRepository, s3Client *s3.Client, bucket string) *DocumentService {
	return &DocumentService{
		repo:     repo,
		s3Client: s3Client,
		bucket:   bucket,
	}
}

func (service *DocumentService) AddDocumentRequest(
	professionalId int,
	clientId int,
	title string,
	description *string,
	dueDate *time.Time,
	status string,
) (int, error) {
	req := models.DocumentRequest{ProfessionalID: professionalId, ClientID: clientId, Title: title, Description: description, DueDate: dueDate, Status: status}
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	id, err := service.repo.AddDocumentRequest(req)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (service *DocumentService) GetDocumentRequestByID(
	id int,
) (models.DocumentRequest, error) {
	req, err := service.repo.GetDocumentRequestByID(id)
	if err != nil {
		return models.DocumentRequest{}, err
	}

	return req, nil
}

func (service *DocumentService) GetDocumentRequestsByProfessional(
	professionalId int,
) ([]models.DocumentRequest, error) {
	reqs, err := service.repo.GetDocumentRequestsByProfessional(professionalId)
	if err != nil {
		return []models.DocumentRequest{}, err
	}

	return reqs, nil
}

func (service *DocumentService) GetDocumentRequestsByClient(
	clientId int,
) ([]models.DocumentRequest, error) {
	reqs, err := service.repo.GetDocumentRequestsByClient(clientId)
	if err != nil {
		return []models.DocumentRequest{}, err
	}

	return reqs, nil
}

func (service *DocumentService) UpdateDocumentRequestStatus(
	id int,
	status string,
) error {
	if err := service.repo.UpdateDocumentRequestStatus(id, status); err != nil {
		return err
	}

	return nil
}

func (service *DocumentService) AddDocumentFile(
	ctx context.Context,
	requestID int,
	fileName string,
	fileSize int64,
	contentType string,
	content io.Reader,
) (int, error) {
	s3Key := fmt.Sprintf("requests/%d/%s", requestID, fileName)

	result, err := service.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(service.bucket),
		Key:         aws.String(s3Key),
		Body:        content,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return 0, errors.ErrInternalServerError{Msg: fmt.Sprintf("Failed to upload to S3: %v", err)}
	}

	fileModel := models.DocumentFile{
		DocumentRequestID: requestID,
		FileName:          fileName,
		FilePath:          s3Key,
		MimeType:          &contentType,
		FileSize:          &fileSize,
		S3VersionID:       result.VersionId,
		UploadedAt:        time.Now(),
	}

	id, err := service.repo.AddDocumentFile(fileModel)
	if err != nil {
		_, cleanupErr := service.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket:    aws.String(service.bucket),
			Key:       aws.String(s3Key),
			VersionId: result.VersionId,
		})

		if cleanupErr != nil {
			return 0, errors.ErrUnprocessableContent{Msg: fmt.Sprintf("Failed to cleanup S3 object after DB failure. Key: %s, Version: %s, Error: %v\n", s3Key, *result.VersionId, cleanupErr)}
		}

		return 0, errors.ErrUnprocessableContent{Msg: fmt.Sprintf("Metadata save failed, file removed from storage: %v", err)}
	}

	return id, nil
}

func (service *DocumentService) GetFilesByRequest() {

}
