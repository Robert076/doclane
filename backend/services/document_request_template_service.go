package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
)

type DocumentRequestTemplateService struct {
	templateRepo        repositories.IDocumentRequestTemplateRepository
	expectedDocTmplRepo repositories.IExpectedDocumentTemplateRepository
	expectedDocRepo     repositories.IExpectedDocumentRepository
	documentRepo        repositories.IDocumentRepository
	userRepo            repositories.IUserRepository
	txManager           repositories.ITxManager
	fileStorage         *FileStorageService
	logger              *slog.Logger
}

func NewDocumentRequestTemplateService(
	templateRepo repositories.IDocumentRequestTemplateRepository,
	expectedDocTmplRepo repositories.IExpectedDocumentTemplateRepository,
	expectedDocRepo repositories.IExpectedDocumentRepository,
	documentRepo repositories.IDocumentRepository,
	userRepo repositories.IUserRepository,
	txManager repositories.ITxManager,
	fileStorage *FileStorageService,
	logger *slog.Logger,
) *DocumentRequestTemplateService {
	return &DocumentRequestTemplateService{
		templateRepo:        templateRepo,
		expectedDocTmplRepo: expectedDocTmplRepo,
		expectedDocRepo:     expectedDocRepo,
		documentRepo:        documentRepo,
		userRepo:            userRepo,
		txManager:           txManager,
		fileStorage:         fileStorage,
		logger:              logger,
	}
}

func (s *DocumentRequestTemplateService) GetTemplatesByProfessionalID(
	ctx context.Context,
	jwtUserID int,
) ([]models.DocumentRequestTemplate, error) {
	templates, err := s.templateRepo.GetDocumentRequestTemplatesByProfessionalID(ctx, jwtUserID)
	if err != nil {
		s.logger.Error("failed to fetch templates",
			slog.Int("professional_id", jwtUserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return templates, nil
}

func (s *DocumentRequestTemplateService) GetTemplateByID(
	ctx context.Context,
	jwtUserID int,
	templateID int,
) (models.DocumentRequestTemplate, error) {
	template, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, templateID)
	if err != nil {
		s.logger.Error("failed to fetch template by id",
			slog.Int("template_id", templateID),
			slog.Any("error", err),
		)
		return models.DocumentRequestTemplate{}, errors.ErrNotFound{Msg: "Template not found."}
	}

	if template.CreatedBy != jwtUserID {
		s.logger.Warn("unauthorized access to template",
			slog.Int("user_id", jwtUserID),
			slog.Int("template_id", templateID),
		)
		return models.DocumentRequestTemplate{}, errors.ErrForbidden{Msg: "You are not allowed to access this template."}
	}

	return template, nil
}

func (s *DocumentRequestTemplateService) AddTemplate(
	ctx context.Context,
	jwtUserID int,
	template models.DocumentRequestTemplate,
) (int, error) {
	if err := ValidateTemplateInput(template); err != nil {
		s.logger.Warn("template create failed validation",
			slog.Int("professional_id", jwtUserID),
			slog.Any("error", err),
		)
		return 0, err
	}

	template.CreatedBy = jwtUserID
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	id, err := s.templateRepo.AddDocumentRequestTemplate(ctx, template)
	if err != nil {
		s.logger.Error("failed to create template",
			slog.Int("professional_id", jwtUserID),
			slog.Any("error", err),
		)
		return 0, err
	}

	s.logger.Info("template created",
		slog.Int("template_id", id),
		slog.Int("professional_id", jwtUserID),
	)
	return id, nil
}

func (s *DocumentRequestTemplateService) AddExpectedDocumentTemplate(
	ctx context.Context,
	jwtUserID int,
	t models.ExpectedDocumentTemplate,
	exampleFile io.Reader,
	exampleFileName string,
	ExampleMimeType string,
	exampleFileSize int64,
) (int, error) {
	template, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, t.DocumentRequestTemplateID)
	if err != nil {
		return 0, errors.ErrNotFound{Msg: "Template not found."}
	}

	if template.CreatedBy != jwtUserID {
		return 0, errors.ErrForbidden{Msg: "You are not allowed to modify this template."}
	}

	if exampleFile != nil {
		if err := ValidateFileInfo(exampleFileName, exampleFileSize); err != nil {
			return 0, err
		}

		s3Key := generateExampleS3Key(exampleFileName)
		result, err := s.fileStorage.UploadFile(ctx, s3Key, exampleFile, ExampleMimeType)
		if err != nil {
			s.logger.Error("s3 upload failed for example file",
				slog.String("key", s3Key),
				slog.Any("error", err),
			)
			return 0, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to upload example file to S3. %v", err)}
		}

		t.ExampleFilePath = &s3Key
		t.ExampleMimeType = &ExampleMimeType

		id, err := s.expectedDocTmplRepo.Add(ctx, t)
		if err != nil {
			s.logger.Error("failed to add expected document template",
				slog.Int("template_id", t.DocumentRequestTemplateID),
				slog.Any("error", err),
			)

			cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cleanupErr := s.fileStorage.DeleteFile(cleanupCtx, s3Key, result.VersionId)
			if cleanupErr != nil {
				s.logger.Error("failed to cleanup example file after db failure",
					slog.String("key", s3Key),
					slog.Any("error", cleanupErr),
				)
			}
			return 0, err
		}

		s.logger.Info("expected document template added with example",
			slog.Int("expected_document_template_id", id),
			slog.Int("template_id", t.DocumentRequestTemplateID),
		)
		return id, nil
	}

	id, err := s.expectedDocTmplRepo.Add(ctx, t)
	if err != nil {
		s.logger.Error("failed to add expected document template",
			slog.Int("template_id", t.DocumentRequestTemplateID),
			slog.Any("error", err),
		)
		return 0, err
	}

	s.logger.Info("expected document template added",
		slog.Int("expected_document_template_id", id),
		slog.Int("template_id", t.DocumentRequestTemplateID),
	)
	return id, nil
}

func (s *DocumentRequestTemplateService) DeleteExpectedDocumentTemplate(
	ctx context.Context,
	jwtUserID int,
	expectedDocTemplateID int,
	templateID int,
) error {
	template, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, templateID)
	if err != nil {
		return errors.ErrNotFound{Msg: "Template not found."}
	}

	if template.CreatedBy != jwtUserID {
		return errors.ErrForbidden{Msg: "You are not allowed to modify this template."}
	}

	if err := s.expectedDocTmplRepo.DeleteByID(ctx, expectedDocTemplateID); err != nil {
		s.logger.Error("failed to delete expected document template",
			slog.Int("expected_document_template_id", expectedDocTemplateID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("expected document template deleted",
		slog.Int("expected_document_template_id", expectedDocTemplateID),
	)
	return nil
}

func (s *DocumentRequestTemplateService) InstantiateTemplate(
	ctx context.Context,
	jwtUserID int,
	templateID int,
	clientID int,
	isScheduled bool,
	scheduledFor *string,
	dueDate *time.Time,
) (int, error) {
	template, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, templateID)
	if err != nil {
		return 0, errors.ErrNotFound{Msg: "Template not found."}
	}

	if template.CreatedBy != jwtUserID {
		return 0, errors.ErrForbidden{Msg: "You are not allowed to instantiate this template."}
	}

	client, err := s.userRepo.GetUserByID(ctx, clientID)
	if err != nil {
		return 0, errors.ErrNotFound{Msg: "Client not found."}
	}

	if client.ProfessionalID == nil || *client.ProfessionalID != jwtUserID {
		s.logger.Warn("unauthorized attempt to instantiate template for unassigned client",
			slog.Int("professional_id", jwtUserID),
			slog.Int("client_id", clientID),
		)
		return 0, errors.ErrForbidden{Msg: "This client is not assigned to you."}
	}

	expectedDocTemplates, err := s.expectedDocTmplRepo.GetByTemplateID(ctx, templateID)
	if err != nil {
		s.logger.Error("failed to fetch expected document templates",
			slog.Int("template_id", templateID),
			slog.Any("error", err),
		)
		return 0, err
	}

	nextDueAt := ComputeNextDueAt(dueDate, template.RecurrenceCron)

	req := models.DocumentRequest{
		ProfessionalID: jwtUserID,
		TemplateID:     &templateID,
		DocumentRequestBase: models.DocumentRequestBase{
			ClientID:       clientID,
			Title:          template.Title,
			Description:    template.Description,
			IsRecurring:    template.IsRecurring,
			RecurrenceCron: template.RecurrenceCron,
			IsScheduled:    isScheduled,
			ScheduledFor:   scheduledFor,
			DueDate:        dueDate,
			NextDueAt:      nextDueAt,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var id int
	err = s.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		var txErr error
		id, txErr = s.documentRepo.AddDocumentRequestWithTx(ctx, req, tx)
		if txErr != nil {
			return txErr
		}

		for _, edt := range expectedDocTemplates {
			ed := models.ExpectedDocument{
				DocumentRequestID: id,
				Title:             edt.Title,
				Description:       edt.Description,
				Status:            "pending",
				ExampleFilePath:   edt.ExampleFilePath,
				ExampleMimeType:   edt.ExampleMimeType,
			}
			if txErr = s.expectedDocRepo.AddExpectedDocumentToRequestWithTx(ctx, tx, ed); txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Error("failed to instantiate template",
			slog.Int("template_id", templateID),
			slog.Int("professional_id", jwtUserID),
			slog.Int("client_id", clientID),
			slog.Any("error", err),
		)
		return 0, err
	}

	s.logger.Info("template instantiated successfully",
		slog.Int("request_id", id),
		slog.Int("template_id", templateID),
		slog.Int("expected_documents", len(expectedDocTemplates)),
	)
	return id, nil
}

func (s *DocumentRequestTemplateService) GetExpectedDocumentTemplatesByTemplateID(ctx context.Context, jwtUserID int, documentRequestTemplateID int) ([]models.ExpectedDocumentTemplate, error) {
	requestTemplate, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, documentRequestTemplateID)
	if err != nil {
		s.logger.Error("failed to retrieve template by id when trying to retrieve expected document templates from it",
			slog.Int("template_id", documentRequestTemplateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return nil, err
	}

	if requestTemplate.CreatedBy != jwtUserID {
		s.logger.Warn("unauthorized access attempted for getting expected document templates for a request template",
			slog.Int("template_id", documentRequestTemplateID),
			slog.Int("user_id", jwtUserID),
		)

		return nil, errors.ErrForbidden{Msg: "This template does not belong to you."}
	}

	documentTemplates, err := s.expectedDocTmplRepo.GetByTemplateID(ctx, documentRequestTemplateID)
	if err != nil {
		s.logger.Error("failed to retrieve document templates by template id",
			slog.Int("template_id", documentRequestTemplateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return nil, err
	}

	return documentTemplates, nil
}

func (s *DocumentRequestTemplateService) PresignExample(ctx context.Context, jwtUserID int, templateID int, expectedDocID int) (string, error) {
	template, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, templateID)
	if err != nil {
		s.logger.Error("failed to retrieve template by id when trying to presign url for example",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return "", err
	}

	if template.CreatedBy != jwtUserID {
		s.logger.Warn("unauthorized access attempted for getting expected document templates for a request template",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
		)

		return "", errors.ErrForbidden{Msg: "This template does not belong to you."}
	}

	expectedDoc, err := s.expectedDocTmplRepo.GetByDocumentID(ctx, expectedDocID)
	if err != nil {
		s.logger.Error("failed to retrieve document template by id when trying to presign url for example",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return "", err
	}

	if expectedDoc.ExampleFilePath == nil {
		s.logger.Warn("attempted document example presign when example does not exist",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return "", errors.ErrBadRequest{Msg: "This template document does not have an example."}
	}

	presignedURL, err := s.fileStorage.GeneratePresignedURL(ctx, *expectedDoc.ExampleFilePath, nil, 15*time.Minute)
	if err != nil {
		s.logger.Error("s3 presign failed for example file",
			slog.Int("expected_doc_id", expectedDocID),
			slog.Any("error", err),
		)
		return "", err
	}
	return presignedURL, nil
}

func (s *DocumentRequestTemplateService) CloseTemplate(ctx context.Context, jwtUserID int, templateID int) error {
	template, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, templateID)
	if err != nil {
		s.logger.Error("failed to retrieve template by id when trying to close it",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return err
	}

	if template.CreatedBy != jwtUserID {
		s.logger.Warn("unauthorized attempt to archive template",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return errors.ErrForbidden{Msg: "You are not allowed to close this template."}
	}

	err = s.templateRepo.CloseDocumentRequestTemplate(ctx, templateID)

	return err
}

func (s *DocumentRequestTemplateService) ReopenTemplate(ctx context.Context, jwtUserID int, templateID int) error {
	template, err := s.templateRepo.GetDocumentRequestTemplateByID(ctx, templateID)
	if err != nil {
		s.logger.Error("failed to retrieve template by id when trying to reopen it",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return err
	}

	if template.CreatedBy != jwtUserID {
		s.logger.Warn("unauthorized attempt to unarchive template",
			slog.Int("template_id", templateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return errors.ErrForbidden{Msg: "You are not allowed to reopen this template."}
	}

	err = s.templateRepo.ReopenDocumentRequestTemplate(ctx, templateID)

	return err
}
