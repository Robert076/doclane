package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Robert076/doclane/backend/events"
	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

func citizenClaims(userID int) types.CallerContext {
	return types.CallerContext{UserID: userID, Role: "citizen"}
}

func configuredUser(id int) models.User {
	phone, street, locality := "0712345678", "Str. Exemplu 1", "Cluj-Napoca"
	return models.User{ID: id, Phone: &phone, Street: &street, Locality: &locality}
}

func newEventBus() *events.EventBus {
	return events.NewEventBus(discardLogger())
}

func TestAddRequest_AdminIsForbidden(t *testing.T) {
	s := &RequestService{logger: discardLogger(), bus: newEventBus()}

	_, err := s.AddRequest(context.Background(), adminClaims(), models.RequestDTOCreate{TemplateID: 1})

	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden for an admin, got %T (%v)", err, err)
	}
}

func TestAddRequest_DepartmentMemberIsForbidden(t *testing.T) {
	s := &RequestService{logger: discardLogger(), bus: newEventBus()}

	_, err := s.AddRequest(context.Background(), memberClaims(2, 5), models.RequestDTOCreate{TemplateID: 1})

	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden for a department member, got %T (%v)", err, err)
	}
}

func TestAddRequest_IncompleteProfileIsRejected(t *testing.T) {
	userRepo := &fakeUserRepo{
		getUserByID: func(ctx context.Context, id int) (models.User, error) {
			return models.User{ID: id}, nil // no phone/street/locality
		},
	}
	s := &RequestService{userRepo: userRepo, logger: discardLogger(), bus: newEventBus()}

	_, err := s.AddRequest(context.Background(), citizenClaims(2), models.RequestDTOCreate{TemplateID: 1})

	if err == nil {
		t.Fatal("expected an error when the citizen profile is not configured")
	}
	if !errors.IsBadRequest(err) {
		t.Errorf("expected ErrBadRequest, got %T (%v)", err, err)
	}
}

func TestAddRequest_MissingTemplateIDFailsValidation(t *testing.T) {
	userRepo := &fakeUserRepo{
		getUserByID: func(ctx context.Context, id int) (models.User, error) {
			return configuredUser(id), nil
		},
	}
	s := &RequestService{userRepo: userRepo, logger: discardLogger(), bus: newEventBus()}

	_, err := s.AddRequest(context.Background(), citizenClaims(2), models.RequestDTOCreate{TemplateID: 0})

	if !errors.IsBadRequest(err) {
		t.Errorf("expected ErrBadRequest for a missing template, got %T (%v)", err, err)
	}
}

func TestAddRequest_UnknownTemplateIsNotFound(t *testing.T) {
	userRepo := &fakeUserRepo{
		getUserByID: func(ctx context.Context, id int) (models.User, error) {
			return configuredUser(id), nil
		},
	}
	templateRepo := &fakeTemplateRepo{
		getRequestTemplateByID: func(ctx context.Context, id int) (models.RequestTemplateDTORead, error) {
			return models.RequestTemplateDTORead{}, errors.ErrNotFound{Msg: "nope"}
		},
	}
	s := &RequestService{
		userRepo:     userRepo,
		templateRepo: templateRepo,
		logger:       discardLogger(),
		bus:          newEventBus(),
	}

	_, err := s.AddRequest(context.Background(), citizenClaims(2), models.RequestDTOCreate{TemplateID: 99})

	if !errors.IsNotFound(err) {
		t.Errorf("expected ErrNotFound for an unknown template, got %T (%v)", err, err)
	}
}

func TestAddRequest_HappyPathCreatesRequestWithExpectedDocs(t *testing.T) {
	const templateID = 7
	const departmentID = 3
	const citizenID = 2

	userRepo := &fakeUserRepo{
		getUserByID: func(ctx context.Context, id int) (models.User, error) {
			return configuredUser(id), nil
		},
	}
	templateRepo := &fakeTemplateRepo{
		getRequestTemplateByID: func(ctx context.Context, id int) (models.RequestTemplateDTORead, error) {
			var tmpl models.RequestTemplateDTORead
			tmpl.ID = id
			tmpl.Title = "ID card renewal"
			tmpl.DepartmentID = departmentID
			return tmpl, nil
		},
	}
	expectedTmplRepo := &fakeExpectedDocTmplRepo{
		getByRequestTemplateID: func(ctx context.Context, id int) ([]models.ExpectedDocumentTemplate, error) {
			return []models.ExpectedDocumentTemplate{
				{Title: "Old ID scan"},
				{Title: "Proof of address"},
			}, nil
		},
	}
	requestRepo := &fakeRequestRepoFull{
		addRequestWithTx: func(ctx context.Context, req models.Request, tx *sql.Tx) (int, error) {
			return 123, nil
		},
	}
	expectedDocRepo := &fakeExpectedDocRepo{}

	s := &RequestService{
		requestRepo:         requestRepo,
		userRepo:            userRepo,
		templateRepo:        templateRepo,
		expectedDocRepo:     expectedDocRepo,
		expectedDocTmplRepo: expectedTmplRepo,
		txManager:           &fakeTxManager{},
		logger:              discardLogger(),
		bus:                 newEventBus(),
	}

	id, err := s.AddRequest(context.Background(), citizenClaims(citizenID),
		models.RequestDTOCreate{TemplateID: templateID})
	if err != nil {
		t.Fatalf("expected no error on the happy path, got %v", err)
	}
	if id == nil || *id != 123 {
		t.Fatalf("expected the created request id 123, got %v", id)
	}
	if len(expectedDocRepo.added) != 2 {
		t.Fatalf("expected 2 expected documents, got %d", len(expectedDocRepo.added))
	}
	for _, ed := range expectedDocRepo.added {
		if ed.Status != "pending" {
			t.Errorf("expected new documents to start as pending, got %q", ed.Status)
		}
	}
}

func TestClaimRequest_CitizenIsForbidden(t *testing.T) {
	s := &RequestService{logger: discardLogger(), bus: newEventBus()}

	err := s.ClaimRequest(context.Background(), citizenClaims(2), 1)

	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden for a citizen, got %T (%v)", err, err)
	}
}

func TestClaimRequest_ClosedRequestRejected(t *testing.T) {
	requestRepo := &fakeRequestRepoFull{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			req := openRequestInDepartment(5, 99)
			req.IsClosed = true
			return req, nil
		},
	}
	s := &RequestService{requestRepo: requestRepo, logger: discardLogger(), bus: newEventBus()}

	err := s.ClaimRequest(context.Background(), memberClaims(2, 5), 1)

	if !errors.IsBadRequest(err) {
		t.Errorf("expected ErrBadRequest for a closed request, got %T (%v)", err, err)
	}
}

func TestClaimRequest_AlreadyClaimedByAnotherIsConflict(t *testing.T) {
	other := 999
	requestRepo := &fakeRequestRepoFull{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			req := openRequestInDepartment(5, 99)
			req.ClaimedBy = &other
			return req, nil
		},
	}
	s := &RequestService{requestRepo: requestRepo, logger: discardLogger(), bus: newEventBus()}

	err := s.ClaimRequest(context.Background(), memberClaims(2, 5), 1)

	if !errors.IsConflict(err) {
		t.Errorf("expected ErrConflict when already claimed by another, got %T (%v)", err, err)
	}
}

func TestClaimRequest_WrongDepartmentForbidden(t *testing.T) {
	requestRepo := &fakeRequestRepoFull{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 99), nil // request in department 5
		},
	}
	s := &RequestService{requestRepo: requestRepo, logger: discardLogger(), bus: newEventBus()}

	err := s.ClaimRequest(context.Background(), memberClaims(2, 7), 1) // member in department 7

	if !errors.IsForbidden(err) {
		t.Errorf("expected ErrForbidden for a different department, got %T (%v)", err, err)
	}
}

func TestClaimRequest_HappyPathClaims(t *testing.T) {
	requestRepo := &fakeRequestRepoFull{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 99), nil
		},
	}
	s := &RequestService{requestRepo: requestRepo, logger: discardLogger(), bus: newEventBus()}

	err := s.ClaimRequest(context.Background(), memberClaims(2, 5), 1)

	if err != nil {
		t.Fatalf("expected the claim to succeed, got %v", err)
	}
	if requestRepo.claimedRequestID != 1 || requestRepo.claimedByUser != 2 {
		t.Errorf("expected request 1 to be claimed by user 2, got request %d by user %d",
			requestRepo.claimedRequestID, requestRepo.claimedByUser)
	}
}

func TestCancelRequest_NonPendingRejected(t *testing.T) {
	requestRepo := &fakeRequestRepoFull{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			req := openRequestInDepartment(5, 2)
			req.ExpectedDocuments = []models.ExpectedDocument{{Status: "uploaded"}}
			return req, nil
		},
	}
	s := &RequestService{requestRepo: requestRepo, logger: discardLogger(), bus: newEventBus()}

	err := s.CancelRequest(context.Background(), memberClaims(2, 5), 1)

	if !errors.IsBadRequest(err) {
		t.Errorf("expected ErrBadRequest when cancelling a non-pending request, got %T (%v)", err, err)
	}
}

func TestCancelRequest_PendingByParticipantSucceeds(t *testing.T) {
	requestRepo := &fakeRequestRepoFull{
		getRequestByID: func(ctx context.Context, id int) (models.RequestDTORead, error) {
			return openRequestInDepartment(5, 2), nil
		},
	}
	s := &RequestService{requestRepo: requestRepo, logger: discardLogger(), bus: newEventBus()}

	err := s.CancelRequest(context.Background(), adminClaims(), 1)

	if err != nil {
		t.Fatalf("expected cancel to succeed for a pending request, got %v", err)
	}
	if requestRepo.cancelledID != 1 {
		t.Errorf("expected request 1 to be cancelled, got %d", requestRepo.cancelledID)
	}
}
