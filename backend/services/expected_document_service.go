package services

import (
	"context"
	"log/slog"

	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
)

type ExpectedRequestService struct {
	expectedDocRepo repositories.IExpectedDocumentRepo
	logger          *slog.Logger
}

func NewExpectedRequestService(expectedDocRepo repositories.IExpectedDocumentRepo, logger *slog.Logger) *ExpectedRequestService {
	return &ExpectedRequestService{
		expectedDocRepo: expectedDocRepo,
		logger:          logger,
	}
}

func (service *ExpectedRequestService) UpdateExpectedDocumentStatus(
	ctx context.Context,
	expectedDocID int,
	status string,
	rejectionReason *string,
) error {
	if status == "rejected" && (rejectionReason == nil || *rejectionReason == "") {
		service.logger.Warn("rejection attempted without a reason", "expectedDocID", expectedDocID)
		return errors.ErrBadRequest{Msg: "Must provide a reason for rejecting the document."}
	}

	if err := service.expectedDocRepo.UpdateExpectedDocumentStatus(ctx, expectedDocID, status, rejectionReason); err != nil {
		service.logger.Error("failed to update expected document status", "expectedDocID", expectedDocID, "status", status, "error", err)
		return err
	}

	service.logger.Info("expected document status updated", "expectedDocID", expectedDocID, "status", status)
	return nil
}
