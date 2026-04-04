package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/robfig/cron/v3"
)

type uploadedFile struct {
	index     int
	s3Key     string
	versionID string
	mimeType  string
}

func validateRequestTemplatePatchDTO(dto models.RequestTemplateDTOPatch) error {
	if dto.Title != nil {
		if strings.TrimSpace(*dto.Title) == "" {
			return errors.ErrBadRequest{Msg: "Title cannot be empty."}
		}
		if len(*dto.Title) > 255 {
			return errors.ErrBadRequest{Msg: "Title cannot exceed 255 characters."}
		}
	}

	if dto.Description != nil && len(*dto.Description) > 1000 {
		return errors.ErrBadRequest{Msg: "Description cannot exceed 1000 characters."}
	}

	if dto.IsRecurring != nil && *dto.IsRecurring && (dto.RecurrenceCron == nil || strings.TrimSpace(*dto.RecurrenceCron) == "") {
		return errors.ErrBadRequest{Msg: "Recurrence cron is required when is_recurring is true."}
	}

	if dto.RecurrenceCron != nil && strings.TrimSpace(*dto.RecurrenceCron) != "" {
		if _, err := cron.ParseStandard(*dto.RecurrenceCron); err != nil {
			return errors.ErrBadRequest{Msg: "Invalid cron expression."}
		}
	}

	return nil
}

func (s *RequestTemplateService) uploadExampleFiles(
	ctx context.Context,
	docs []types.ExpectedDocumentTemplateInput,
) (map[int]uploadedFile, func(), error) {
	var uploads []uploadedFile

	rollbackS3 := func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, u := range uploads {
			var vID *string
			if u.versionID != "" {
				vID = &u.versionID
			}
			if err := s.fileStorage.DeleteFile(cleanupCtx, u.s3Key, vID); err != nil {
				s.logger.Error("s3 cleanup failed during rollback",
					slog.String("key", u.s3Key),
					slog.Any("error", err),
				)
			}
		}
	}

	for i, doc := range docs {
		if doc.ExampleFile == nil {
			continue
		}
		if err := ValidateFileInfo(doc.ExampleFileName, doc.ExampleFileSize); err != nil {
			rollbackS3()
			return nil, nil, err
		}

		s3Key := s.fileStorage.GenerateS3Key(doc.ExampleFileName, "examples")
		result, err := s.fileStorage.UploadFile(ctx, s3Key, doc.ExampleFile, doc.ExampleMimeType)
		if err != nil {
			rollbackS3()
			return nil, nil, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to upload example file %d to S3: %v", i, err)}
		}

		versionID := ""
		if result.VersionId != nil {
			versionID = *result.VersionId
		}

		uploads = append(uploads, uploadedFile{
			index:     i,
			s3Key:     s3Key,
			versionID: versionID,
			mimeType:  doc.ExampleMimeType,
		})
	}

	uploadByIndex := make(map[int]uploadedFile, len(uploads))
	for _, u := range uploads {
		uploadByIndex[u.index] = u
	}

	return uploadByIndex, rollbackS3, nil
}

func (s *RequestTemplateService) insertTemplateWithDocsTx(
	ctx context.Context,
	template models.RequestTemplate,
	docs []types.ExpectedDocumentTemplateInput,
	uploadByIndex map[int]uploadedFile,
) (int, error) {
	var templateID int
	err := s.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		var txErr error
		templateID, txErr = s.templateRepo.AddRequestTemplateWithTx(ctx, tx, template)
		if txErr != nil {
			return txErr
		}

		for i, doc := range docs {
			ed := models.ExpectedDocumentTemplate{
				RequestTemplateID: templateID,
				Title:             doc.Title,
				Description:       doc.Description,
			}
			if u, ok := uploadByIndex[i]; ok {
				ed.ExampleFilePath = &u.s3Key
				ed.ExampleMimeType = &u.mimeType
			}
			if _, txErr = s.expectedDocTmplRepo.AddWithTx(ctx, tx, ed); txErr != nil {
				return txErr
			}
		}
		return nil
	})
	return templateID, err
}

func (s *RequestTemplateService) checkUserCanAccessTemplate(ctx context.Context, claims types.JWTClaims, requestTemplateID int) (*models.RequestTemplateDTORead, error) {
	template, err := s.templateRepo.GetRequestTemplateByID(ctx, requestTemplateID)
	if err != nil {
		s.logger.Error("error getting template from db",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("template_id", requestTemplateID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if claims.IsAdmin() {
		return &template, nil
	}

	isDepartmentMatch := claims.DepartmentID != nil && *claims.DepartmentID == template.DepartmentID
	if !isDepartmentMatch {
		s.logger.Warn("unauthorized access attempted for template",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("template_id", requestTemplateID),
		)
		return nil, errors.ErrForbidden{Msg: "You don't have access to this template."}
	}

	return &template, nil
}
