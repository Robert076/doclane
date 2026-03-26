package services

import (
	"context"
	"log/slog"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
)

type ExpectedDocumentService struct {
	expectedDocRepo repositories.IExpectedDocumentRepo
	requestRepo     repositories.IRequestRepo
	logger          *slog.Logger
}

func NewExpectedDocumentService(expectedDocRepo repositories.IExpectedDocumentRepo, requestRepo repositories.IRequestRepo, logger *slog.Logger) *ExpectedDocumentService {
	return &ExpectedDocumentService{
		expectedDocRepo: expectedDocRepo,
		requestRepo:     requestRepo,
		logger:          logger,
	}
}

func (service *ExpectedDocumentService) UpdateExpectedDocumentStatus(
	ctx context.Context,
	jwtUserID int,
	expectedDocID int,
	status string,
	rejectionReason *string,
) (*models.ExpectedDocument, error) {
	if status == "rejected" && (rejectionReason == nil || *rejectionReason == "") {
		service.logger.Warn("rejection attempted without a reason", "expectedDocID", expectedDocID)
		return nil, errors.ErrBadRequest{Msg: "Must provide a reason for rejecting the document."}
	}

	expectedDoc, err := service.expectedDocRepo.GetExpectedDocumentByID(ctx, expectedDocID)
	if err != nil {
		return nil, errors.ErrNotFound{Msg: "Expected document not found."}
	}

	req, err := service.requestRepo.GetRequestByID(ctx, expectedDoc.RequestID)
	if err != nil {
		return nil, errors.ErrNotFound{Msg: "Request not found."}
	}

	if req.ProfessionalID != jwtUserID {
		service.logger.Warn("unauthorized status update attempt",
			slog.Int("user_id", jwtUserID),
			slog.Int("expected_doc_id", expectedDocID),
		)
		return nil, errors.ErrForbidden{Msg: "Only the professional can update document status."}
	}

	if err := service.expectedDocRepo.UpdateExpectedDocumentStatus(ctx, expectedDocID, status, rejectionReason); err != nil {
		service.logger.Error("failed to update expected document status",
			slog.Int("expected_doc_id", expectedDocID),
			slog.String("status", status),
			slog.Any("error", err),
		)
		return nil, err
	}

	updated, err := service.expectedDocRepo.GetExpectedDocumentByID(ctx, expectedDocID)
	if err != nil {
		return nil, err
	}

	service.logger.Info("expected document status updated",
		slog.Int("expected_doc_id", expectedDocID),
		slog.String("status", status),
	)
	return &updated, nil
}
