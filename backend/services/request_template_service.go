package services

import (
	"context"
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
	txManager           repositories.ITxManager
	fileStorage         IFileStorageService
	logger              *slog.Logger
}

func NewRequestTemplateService(
	templateRepo repositories.IRequestTemplateRepo,
	expectedDocTmplRepo repositories.IExpectedDocumentTemplateRepo,
	expectedDocRepo repositories.IExpectedDocumentRepo,
	documentRepo repositories.IRequestRepo,
	txManager repositories.ITxManager,
	fileStorage IFileStorageService,
	logger *slog.Logger,
) *RequestTemplateService {
	return &RequestTemplateService{
		templateRepo:        templateRepo,
		expectedDocTmplRepo: expectedDocTmplRepo,
		expectedDocRepo:     expectedDocRepo,
		documentRepo:        documentRepo,
		txManager:           txManager,
		fileStorage:         fileStorage,
		logger:              logger,
	}
}

func (s *RequestTemplateService) GetRequestTemplates(ctx context.Context, claims types.JWTClaims) ([]models.RequestTemplateDTORead, error) {
	if claims.IsAdmin() {
		templates, err := s.templateRepo.GetAllRequestTemplates(ctx)
		if err != nil {
			s.logger.Error("failed to fetch all templates",
				slog.Int("jwt_user_id", claims.UserID),
				slog.Any("error", err),
			)
			return nil, err
		}
		s.logger.Info("fetched all templates successfully",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return templates, nil
	}

	if claims.DepartmentID == nil {
		return nil, errors.ErrForbidden{Msg: "You are not part of a department."}
	}

	templates, err := s.templateRepo.GetRequestTemplatesByDepartment(ctx, *claims.DepartmentID)
	if err != nil {
		s.logger.Error("failed to fetch templates by department",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("department_id", *claims.DepartmentID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("fetched templates by department successfully",
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("department_id", *claims.DepartmentID),
	)
	return templates, nil
}

func (s *RequestTemplateService) GetRequestTemplateByID(ctx context.Context, claims types.JWTClaims, requestTemplateID int) (*models.RequestTemplate, error) {
	return s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID)
}

func (s *RequestTemplateService) AddRequestTemplateWithDocuments(
	ctx context.Context,
	claims types.JWTClaims,
	template models.RequestTemplate,
	docs []types.ExpectedDocumentTemplateInput,
) (*int, error) {
	if !claims.IsAdmin() && !claims.IsDepartmentMember() {
		s.logger.Warn("unauthorized attempt to create template",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "You are not allowed to create templates."}
	}

	if err := ValidateRequestTemplateInput(template); err != nil {
		return nil, err
	}

	uploadByIndex, rollbackS3, err := s.uploadExampleFiles(ctx, docs)
	if err != nil {
		return nil, err
	}

	template.CreatedBy = claims.UserID

	templateID, err := s.insertTemplateWithDocsTx(ctx, template, docs, uploadByIndex)
	if err != nil {
		rollbackS3()
		s.logger.Error("transaction failed, rolled back S3 uploads",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("template created with documents",
		slog.Int("template_id", templateID),
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("document_count", len(docs)),
	)
	return &templateID, nil
}

func (s *RequestTemplateService) DeleteExpectedDocumentTemplate(
	ctx context.Context,
	claims types.JWTClaims,
	expectedDocRequestTemplateID int,
	requestTemplateID int,
) error {
	if _, err := s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID); err != nil {
		return err
	}

	if err := s.expectedDocTmplRepo.DeleteByID(ctx, expectedDocRequestTemplateID); err != nil {
		s.logger.Error("failed to delete expected document template",
			slog.Int("expected_document_template_id", expectedDocRequestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("expected document template deleted",
		slog.Int("expected_document_template_id", expectedDocRequestTemplateID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *RequestTemplateService) GetExpectedDocumentTemplatesByRequestTemplateID(ctx context.Context, claims types.JWTClaims, requestTemplateID int) ([]models.ExpectedDocumentTemplate, error) {
	if _, err := s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID); err != nil {
		return nil, err
	}

	documentRequestTemplates, err := s.expectedDocTmplRepo.GetByRequestTemplateID(ctx, requestTemplateID)
	if err != nil {
		s.logger.Error("failed to retrieve document templates by template id",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	return documentRequestTemplates, nil
}

func (s *RequestTemplateService) PresignExample(ctx context.Context, claims types.JWTClaims, requestTemplateID int, expectedDocID int) (*string, error) {
	template, err := s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID)
	if err != nil {
		return nil, err
	}

	expectedDoc, err := s.expectedDocTmplRepo.GetByDocumentID(ctx, expectedDocID)
	if err != nil {
		s.logger.Error("failed to retrieve document template by id when trying to presign url for example",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if expectedDoc.RequestTemplateID != template.ID {
		s.logger.Warn("unauthorized retrieval attempt for example document when presigning",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("example_document_id", expectedDocID),
		)
		return nil, errors.ErrForbidden{Msg: "You are not allowed to view this file."}
	}

	if expectedDoc.ExampleFilePath == nil {
		s.logger.Warn("attempted document example presign when example does not exist",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrBadRequest{Msg: "This template document does not have an example."}
	}

	presignedURL, err := s.fileStorage.GeneratePresignedURL(ctx, *expectedDoc.ExampleFilePath, nil, 15*time.Minute)
	if err != nil {
		s.logger.Error("s3 presign failed for example file",
			slog.Int("expected_doc_id", expectedDocID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return &presignedURL, nil
}

func (s *RequestTemplateService) CloseRequestTemplate(ctx context.Context, claims types.JWTClaims, requestTemplateID int) error {
	if _, err := s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID); err != nil {
		return err
	}

	if err := s.templateRepo.CloseRequestTemplate(ctx, requestTemplateID); err != nil {
		s.logger.Error("failed to close request template",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("request template closed successfully",
		slog.Int("template_id", requestTemplateID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *RequestTemplateService) ReopenRequestTemplate(ctx context.Context, claims types.JWTClaims, requestTemplateID int) error {
	if _, err := s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID); err != nil {
		return err
	}

	if err := s.templateRepo.ReopenRequestTemplate(ctx, requestTemplateID); err != nil {
		s.logger.Error("failed to reopen request template",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("request template reopened successfully",
		slog.Int("template_id", requestTemplateID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *RequestTemplateService) DeleteRequestTemplate(ctx context.Context, claims types.JWTClaims, requestTemplateID int) error {
	if _, err := s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID); err != nil {
		return err
	}

	if err := s.templateRepo.DeleteRequestTemplate(ctx, requestTemplateID); err != nil {
		s.logger.Error("failed to delete request template",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("request template deleted successfully",
		slog.Int("template_id", requestTemplateID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *RequestTemplateService) PatchRequestTemplate(ctx context.Context, claims types.JWTClaims, requestTemplateID int, dto models.RequestTemplateDTOPatch) error {
	if err := validateRequestTemplatePatchDTO(dto); err != nil {
		return err
	}

	if _, err := s.checkUserCanAccessTemplate(ctx, claims, requestTemplateID); err != nil {
		return err
	}

	if err := s.templateRepo.PatchRequestTemplate(ctx, requestTemplateID, dto); err != nil {
		s.logger.Error("failed to patch request template",
			slog.Int("template_id", requestTemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("request template patched successfully",
		slog.Int("template_id", requestTemplateID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}
