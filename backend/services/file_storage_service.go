package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type IFileStorageService interface {
	UploadFile(ctx context.Context, key string, content io.Reader, contentType string) (*s3.PutObjectOutput, error)
	DeleteFile(ctx context.Context, key string, versionID *string) error
	GeneratePresignedURL(ctx context.Context, key string, versionID *string, expiry time.Duration) (string, error)
	GenerateS3Key(fileName string, prefix string) string
}

type FileStorageService struct {
	s3Client *s3.Client
	bucket   string
	logger   *slog.Logger
}

func NewFileStorageService(s3Client *s3.Client, bucket string, logger *slog.Logger) *FileStorageService {
	return &FileStorageService{
		s3Client: s3Client,
		bucket:   bucket,
		logger:   logger,
	}
}

func (s *FileStorageService) UploadFile(ctx context.Context, key string, content io.Reader, contentType string) (*s3.PutObjectOutput, error) {
	result, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        content,
		ContentType: aws.String(contentType),
	})

	return result, err
}
func (s *FileStorageService) DeleteFile(ctx context.Context, key string, versionID *string) error {
	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket:    aws.String(s.bucket),
		Key:       aws.String(key),
		VersionId: versionID,
	})

	return err
}

func (s *FileStorageService) GeneratePresignedURL(ctx context.Context, key string, versionID *string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.s3Client)
	presignParams := &s3.GetObjectInput{
		Bucket:    aws.String(s.bucket),
		Key:       aws.String(key),
		VersionId: versionID,
	}

	presignedReq, err := presignClient.PresignGetObject(ctx, presignParams, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	return presignedReq.URL, err
}

func (s *FileStorageService) GenerateS3Key(fileName string, prefix string) string {
	cleanFileName := filepath.Base(fileName)
	uniqueID := uuid.New().String()
	return fmt.Sprintf("%s/%s-%s", prefix, uniqueID, cleanFileName)
}
