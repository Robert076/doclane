package services

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

type RequestTemplateService struct {
	templateRepo        repositories.IRequestTemplateRepo
	expectedDocTmplRepo repositories.IExpectedDocumentTemplateRepo
	expectedDocRepo     repositories.IExpectedDocumentRepo
	documentRepo        repositories.IRequestRepo
	userRepo            repositories.IUserRepo
	txManager           repositories.ITxManager
	fileStorage         IFileStorageService
	logger              *slog.Logger
}

func NewRequestTemplateService(
	templateRepo repositories.IRequestTemplateRepo,
	expectedDocTmplRepo repositories.IExpectedDocumentTemplateRepo,
	expectedDocRepo repositories.IExpectedDocumentRepo,
	documentRepo repositories.IRequestRepo,
	userRepo repositories.IUserRepo,
	txManager repositories.ITxManager,
	fileStorage IFileStorageService,
	logger *slog.Logger,
) *RequestTemplateService {
	return &RequestTemplateService{
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

func (s *RequestTemplateService) GetRequestTemplatesByProfessionalID(
	ctx context.Context,
	jwtUserID int,
) ([]models.RequestTemplate, error) {
	templates, err := s.templateRepo.GetRequestTemplatesByProfessionalID(ctx, jwtUserID)
	if err != nil {
		s.logger.Error("failed to fetch templates",
			slog.Int("professional_id", jwtUserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return templates, nil
}

func (s *RequestTemplateService) GetRequestTemplateByID(
	ctx context.Context,
	jwtUserID int,
	requestTemplateID int,
) (*models.RequestTemplate, error) {
	return s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID)
}

func (s *RequestTemplateService) AddRequestTemplateWithDocuments(
	ctx context.Context,
	jwtUserID int,
	template models.RequestTemplate,
	docs []types.ExpectedDocumentTemplateInput,
) (*int, error) {
	if err := ValidateRequestTemplateInput(template); err != nil {
		return nil, err
	}

	uploadByIndex, rollbackS3, err := s.uploadExampleFiles(ctx, docs)
	if err != nil {
		return nil, err
	}

	template.CreatedBy = jwtUserID
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	templateID, err := s.insertTemplateWithDocsTx(ctx, template, docs, uploadByIndex)
	if err != nil {
		rollbackS3()
		s.logger.Error("transaction failed, rolled back S3 uploads",
			slog.Int("professional_id", jwtUserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("template created with documents",
		slog.Int("template_id", templateID),
		slog.Int("document_count", len(docs)),
	)
	return &templateID, nil
}

func (s *RequestTemplateService) DeleteExpectedDocumentTemplate(
	ctx context.Context,
	jwtUserID int,
	expectedDocRequestTemplateID int,
	requestTemplateID int,
) error {
	if _, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID); err != nil {
		return err
	}

	if err := s.expectedDocTmplRepo.DeleteByID(ctx, expectedDocRequestTemplateID); err != nil {
		s.logger.Error("failed to delete expected document template",
			slog.Int("expected_document_template_id", expectedDocRequestTemplateID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("expected document template deleted",
		slog.Int("expected_document_template_id", expectedDocRequestTemplateID),
	)
	return nil
}

func (s *RequestTemplateService) InstantiateRequestTemplate(
	ctx context.Context,
	jwtUserID int,
	requestTemplateID int,
	clientID int,
	isScheduled bool,
	scheduledFor *string,
	dueDate *time.Time,
) (*int, error) {
	template, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID)
	if err != nil {
		return nil, err
	}

	client, err := s.userRepo.GetUserByID(ctx, clientID)
	if err != nil {
		return nil, errors.ErrNotFound{Msg: "Client not found."}
	}

	if client.ProfessionalID == nil || *client.ProfessionalID != jwtUserID {
		s.logger.Warn("unauthorized attempt to instantiate template for unassigned client",
			slog.Int("professional_id", jwtUserID),
			slog.Int("client_id", clientID),
		)
		return nil, errors.ErrForbidden{Msg: "This client is not assigned to you."}
	}

	expectedDocRequestTemplates, err := s.expectedDocTmplRepo.GetByRequestTemplateID(ctx, requestTemplateID)
	if err != nil {
		s.logger.Error("failed to fetch expected document templates",
			slog.Int("template_id", requestTemplateID),
			slog.Any("error", err),
		)
		return nil, err
	}

	nextDueAt := ComputeNextDueAt(dueDate, template.RecurrenceCron)

	req := models.Request{
		ProfessionalID:    jwtUserID,
		RequestTemplateID: &requestTemplateID,
		RequestBase: models.RequestBase{
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
		id, txErr = s.documentRepo.AddRequestWithTx(ctx, req, tx)
		if txErr != nil {
			return txErr
		}

		for _, edt := range expectedDocRequestTemplates {
			ed := models.ExpectedDocument{
				RequestID:       id,
				Title:           edt.Title,
				Description:     edt.Description,
				Status:          "pending",
				ExampleFilePath: edt.ExampleFilePath,
				ExampleMimeType: edt.ExampleMimeType,
			}
			if txErr = s.expectedDocRepo.AddExpectedDocumentToRequestWithTx(ctx, tx, ed); txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Error("failed to instantiate template",
			slog.Int("template_id", requestTemplateID),
			slog.Int("professional_id", jwtUserID),
			slog.Int("client_id", clientID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("template instantiated successfully",
		slog.Int("request_id", id),
		slog.Int("template_id", requestTemplateID),
		slog.Int("expected_documents", len(expectedDocRequestTemplates)),
	)
	return &id, nil
}

func (s *RequestTemplateService) GetExpectedDocumentTemplatesByRequestTemplateID(ctx context.Context, jwtUserID int, requestTemplateID int) ([]models.ExpectedDocumentTemplate, error) {
	if _, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID); err != nil {
		return nil, err
	}

	documentRequestTemplates, err := s.expectedDocTmplRepo.GetByRequestTemplateID(ctx, requestTemplateID)
	if err != nil {
		s.logger.Error("failed to retrieve document templates by template id",
			slog.Int("template_id", requestTemplateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return nil, err
	}

	return documentRequestTemplates, nil
}

func (s *RequestTemplateService) PresignExample(ctx context.Context, jwtUserID int, requestTemplateID int, expectedDocID int) (*string, error) {
	template, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID)
	if err != nil {
		return nil, err
	}

	expectedDoc, err := s.expectedDocTmplRepo.GetByDocumentID(ctx, expectedDocID)
	if err != nil {
		s.logger.Error("failed to retrieve document template by id when trying to presign url for example",
			slog.Int("template_id", requestTemplateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return nil, err
	}

	if expectedDoc.RequestTemplateID != template.ID {
		s.logger.Warn("unauthorized retrieval attempt for example document when presigning",
			slog.Int("template_id", requestTemplateID),
			slog.Int("user_id", jwtUserID),
			slog.Int("example_document_id", expectedDocID),
		)
		return nil, errors.ErrForbidden{Msg: "You are not allowed to view this file."}
	}

	if expectedDoc.ExampleFilePath == nil {
		s.logger.Warn("attempted document example presign when example does not exist",
			slog.Int("template_id", requestTemplateID),
			slog.Int("user_id", jwtUserID),
		)

		return nil, errors.ErrBadRequest{Msg: "This template document does not have an example."}
	}

	presignedURL, err := s.fileStorage.GeneratePresignedURL(ctx, *expectedDoc.ExampleFilePath, nil, 15*time.Minute)
	if err != nil {
		s.logger.Error("s3 presign failed for example file",
			slog.Int("expected_doc_id", expectedDocID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return &presignedURL, nil
}

func (s *RequestTemplateService) CloseRequestTemplate(ctx context.Context, jwtUserID int, requestTemplateID int) error {
	if _, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID); err != nil {
		return err
	}

	return s.templateRepo.CloseRequestTemplate(ctx, requestTemplateID)
}

func (s *RequestTemplateService) ReopenRequestTemplate(ctx context.Context, jwtUserID int, requestTemplateID int) error {
	if _, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID); err != nil {
		return err
	}

	return s.templateRepo.ReopenRequestTemplate(ctx, requestTemplateID)
}

func (s *RequestTemplateService) DeleteRequestTemplate(ctx context.Context, jwtUserID int, requestTemplateID int) error {
	if _, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID); err != nil {
		return err
	}

	return s.templateRepo.DeleteRequestTemplate(ctx, requestTemplateID)
}

func (s *RequestTemplateService) PatchRequestTemplate(ctx context.Context, jwtUserID int, requestTemplateID int, dto models.RequestTemplateDTOPatch) error {
	if err := validateRequestTemplatePatchDTO(dto); err != nil {
		return err
	}

	if _, err := s.checkUserOwnsTemplate(ctx, jwtUserID, requestTemplateID); err != nil {
		return err
	}

	return s.templateRepo.PatchRequestTemplate(ctx, requestTemplateID, dto)
}
