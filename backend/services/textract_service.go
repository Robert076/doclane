package services

import (
	"context"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type ITextractService interface {
	ExtractText(ctx context.Context, s3Key string) (string, error)
}

type TextractService struct {
	client *textract.Client
	bucket string
	logger *slog.Logger
}

func NewTextractService(client *textract.Client, bucket string, logger *slog.Logger) *TextractService {
	return &TextractService{
		client: client,
		bucket: bucket,
		logger: logger,
	}
}

func (s *TextractService) ExtractText(ctx context.Context, s3Key string) (string, error) {
	output, err := s.client.DetectDocumentText(ctx, &textract.DetectDocumentTextInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: aws.String(s.bucket),
				Name:   aws.String(s3Key),
			},
		},
	})
	if err != nil {
		s.logger.Error("textract failed",
			slog.String("s3_key", s3Key),
			slog.Any("error", err),
		)
		return "", err
	}

	var lines []string
	for _, block := range output.Blocks {
		if block.BlockType == types.BlockTypeLine && block.Text != nil {
			lines = append(lines, *block.Text)
		}
	}

	return strings.Join(lines, "\n"), nil
}
