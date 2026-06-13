package services

import (
	"testing"
	"time"

	"github.com/Robert076/doclane/backend/models"
)

func TestValidateRequestInput_ValidInputPasses(t *testing.T) {
	dto := models.RequestDTOCreate{TemplateID: 1}

	if err := ValidateRequestInput(dto); err != nil {
		t.Errorf("expected no error for valid input, got %v", err)
	}
}

func TestValidateRequestInput_MissingTemplateFails(t *testing.T) {
	dto := models.RequestDTOCreate{TemplateID: 0}

	if err := ValidateRequestInput(dto); err == nil {
		t.Error("expected an error when template id is 0")
	}
}

func TestValidateRequestInput_PastDueDateFails(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	dto := models.RequestDTOCreate{TemplateID: 1, DueDate: &past}

	if err := ValidateRequestInput(dto); err == nil {
		t.Error("expected an error when due date is in the past")
	}
}

func TestValidateRequestInput_FutureDueDatePasses(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	dto := models.RequestDTOCreate{TemplateID: 1, DueDate: &future}

	if err := ValidateRequestInput(dto); err != nil {
		t.Errorf("expected no error for a future due date, got %v", err)
	}
}

func TestValidateRequestInput_ScheduledWithoutScheduledForFails(t *testing.T) {
	dto := models.RequestDTOCreate{TemplateID: 1, IsScheduled: true, ScheduledFor: nil}

	if err := ValidateRequestInput(dto); err == nil {
		t.Error("expected an error when a scheduled request has no scheduled_for")
	}
}

func TestValidateRequestInput_ScheduledWithScheduledForPasses(t *testing.T) {
	when := time.Now().Add(time.Hour)
	dto := models.RequestDTOCreate{TemplateID: 1, IsScheduled: true, ScheduledFor: &when}

	if err := ValidateRequestInput(dto); err != nil {
		t.Errorf("expected no error when scheduled_for is set, got %v", err)
	}
}

func TestComputeStatus_NoDocumentsAndNoUploadIsPending(t *testing.T) {
	status := ComputeStatus(nil, nil, nil)

	if status != "pending" {
		t.Errorf("expected pending, got %q", status)
	}
}

func TestComputeStatus_AllDocsAcceptedIsUploaded(t *testing.T) {
	docs := []models.ExpectedDocument{
		{Status: "accepted"},
		{Status: "uploaded"},
	}

	status := ComputeStatus(nil, nil, docs)

	if status != "uploaded" {
		t.Errorf("expected uploaded, got %q", status)
	}
}

func TestComputeStatus_PendingDocWithNoDueDateIsPending(t *testing.T) {
	docs := []models.ExpectedDocument{
		{Status: "pending"},
	}

	status := ComputeStatus(nil, nil, docs)

	if status != "pending" {
		t.Errorf("expected pending, got %q", status)
	}
}

func TestComputeStatus_PastDueWithPendingDocIsOverdue(t *testing.T) {
	due := time.Now().Add(-1 * time.Hour)
	docs := []models.ExpectedDocument{
		{Status: "pending"},
	}

	status := ComputeStatus(nil, &due, docs)

	if status != "overdue" {
		t.Errorf("expected overdue, got %q", status)
	}
}

func TestComputeStatus_PastDueButAllUploadedIsUploaded(t *testing.T) {
	due := time.Now().Add(-1 * time.Hour)
	docs := []models.ExpectedDocument{
		{Status: "uploaded"},
	}

	status := ComputeStatus(nil, &due, docs)

	if status != "uploaded" {
		t.Errorf("expected uploaded, got %q", status)
	}
}

func TestComputeStatus_FutureDueWithPendingDocIsPending(t *testing.T) {
	due := time.Now().Add(24 * time.Hour)
	docs := []models.ExpectedDocument{
		{Status: "pending"},
	}

	status := ComputeStatus(nil, &due, docs)

	if status != "pending" {
		t.Errorf("expected pending, got %q", status)
	}
}

func TestValidatePatchDTO_ValidTitlePasses(t *testing.T) {
	if err := ValidatePatchDTO(models.RequestDTOPatch{Title: "Valid title"}); err != nil {
		t.Errorf("expected no error for a valid title, got %v", err)
	}
}

func TestValidatePatchDTO_TooShortTitleFails(t *testing.T) {
	if err := ValidatePatchDTO(models.RequestDTOPatch{Title: "ab"}); err == nil {
		t.Error("expected an error for a title shorter than 3 characters")
	}
}

func TestValidatePatchDTO_TooLongTitleFails(t *testing.T) {
	long := "this title is definitely longer than thirty characters"
	if err := ValidatePatchDTO(models.RequestDTOPatch{Title: long}); err == nil {
		t.Error("expected an error for a title longer than 30 characters")
	}
}

func TestValidateRequestTemplateInput_ValidTemplatePasses(t *testing.T) {
	tmpl := models.RequestTemplate{Title: "Birth certificate request"}

	if err := ValidateRequestTemplateInput(tmpl); err != nil {
		t.Errorf("expected no error for a valid template, got %v", err)
	}
}

func TestValidateRequestTemplateInput_TooShortTitleFails(t *testing.T) {
	tmpl := models.RequestTemplate{Title: "ab"}

	if err := ValidateRequestTemplateInput(tmpl); err == nil {
		t.Error("expected an error for a title shorter than 3 characters")
	}
}

func TestValidateRequestTemplateInput_RecurringWithoutCronFails(t *testing.T) {
	tmpl := models.RequestTemplate{Title: "Monthly report", IsRecurring: true, RecurrenceCron: nil}

	if err := ValidateRequestTemplateInput(tmpl); err == nil {
		t.Error("expected an error for a recurring template without a cron expression")
	}
}

func TestValidateRequestTemplateInput_RecurringWithValidCronPasses(t *testing.T) {
	cron := "0 0 1 * *"
	tmpl := models.RequestTemplate{Title: "Monthly report", IsRecurring: true, RecurrenceCron: &cron}

	if err := ValidateRequestTemplateInput(tmpl); err != nil {
		t.Errorf("expected no error for a valid cron expression, got %v", err)
	}
}

func TestValidateRequestTemplateInput_InvalidCronFails(t *testing.T) {
	cron := "not a cron"
	tmpl := models.RequestTemplate{Title: "Monthly report", IsRecurring: true, RecurrenceCron: &cron}

	if err := ValidateRequestTemplateInput(tmpl); err == nil {
		t.Error("expected an error for an invalid cron expression")
	}
}

func TestComputeNextDueAt_ReturnsDueDateWhenProvided(t *testing.T) {
	due := time.Now().Add(48 * time.Hour)

	got := ComputeNextDueAt(&due, nil)

	if got == nil || !got.Equal(due) {
		t.Errorf("expected the provided due date to be returned, got %v", got)
	}
}

func TestComputeNextDueAt_NilWhenNoDueDateAndNoCron(t *testing.T) {
	if got := ComputeNextDueAt(nil, nil); got != nil {
		t.Errorf("expected nil when neither due date nor cron is provided, got %v", got)
	}
}

func TestComputeNextDueAt_ComputesFromCron(t *testing.T) {
	cron := "0 0 * * *" // every day at midnight
	now := time.Now()

	got := ComputeNextDueAt(nil, &cron)

	if got == nil {
		t.Fatal("expected a computed next due date from the cron expression, got nil")
	}
	if !got.After(now) {
		t.Errorf("expected the next due date to be in the future, got %v", got)
	}
}

func TestComputeNextDueAt_NilForInvalidCron(t *testing.T) {
	cron := "totally invalid"

	if got := ComputeNextDueAt(nil, &cron); got != nil {
		t.Errorf("expected nil for an invalid cron expression, got %v", got)
	}
}

func TestValidateFileInfo_ValidPdfPasses(t *testing.T) {
	if err := ValidateFileInfo("document.pdf", 1024); err != nil {
		t.Errorf("expected no error for a valid pdf, got %v", err)
	}
}

func TestValidateFileInfo_EmptyFileFails(t *testing.T) {
	if err := ValidateFileInfo("document.pdf", 0); err == nil {
		t.Error("expected an error for an empty file")
	}
}

func TestValidateFileInfo_TooLargeFileFails(t *testing.T) {
	tooBig := int64(21 * 1024 * 1024)
	if err := ValidateFileInfo("document.pdf", tooBig); err == nil {
		t.Error("expected an error for a file larger than 20MB")
	}
}

func TestValidateFileInfo_DisallowedExtensionFails(t *testing.T) {
	if err := ValidateFileInfo("malware.exe", 1024); err == nil {
		t.Error("expected an error for a disallowed extension")
	}
}

func TestValidateFileInfo_AllowedImageExtensionPasses(t *testing.T) {
	if err := ValidateFileInfo("scan.jpeg", 2048); err != nil {
		t.Errorf("expected no error for a jpeg, got %v", err)
	}
}

func TestValidateRequestTemplatePatchDTO_EmptyPatchPasses(t *testing.T) {
	if err := validateRequestTemplatePatchDTO(models.RequestTemplateDTOPatch{}); err != nil {
		t.Errorf("expected no error for an empty patch, got %v", err)
	}
}

func TestValidateRequestTemplatePatchDTO_BlankTitleFails(t *testing.T) {
	blank := "   "
	dto := models.RequestTemplateDTOPatch{Title: &blank}

	if err := validateRequestTemplatePatchDTO(dto); err == nil {
		t.Error("expected an error for a blank title")
	}
}

func TestValidateRequestTemplatePatchDTO_RecurringWithoutCronFails(t *testing.T) {
	isRecurring := true
	dto := models.RequestTemplateDTOPatch{IsRecurring: &isRecurring}

	if err := validateRequestTemplatePatchDTO(dto); err == nil {
		t.Error("expected an error when is_recurring is true with no cron")
	}
}

func TestValidateRequestTemplatePatchDTO_InvalidCronFails(t *testing.T) {
	bad := "nope"
	dto := models.RequestTemplateDTOPatch{RecurrenceCron: &bad}

	if err := validateRequestTemplatePatchDTO(dto); err == nil {
		t.Error("expected an error for an invalid cron expression")
	}
}

func TestValidateTagDTO_ValidNameAndColorPasses(t *testing.T) {
	if err := validateTagDTO("Urgent", "#ff5722"); err != nil {
		t.Errorf("expected no error for a valid tag, got %v", err)
	}
}

func TestValidateTagDTO_EmptyColorPasses(t *testing.T) {
	if err := validateTagDTO("Urgent", ""); err != nil {
		t.Errorf("expected no error when color is empty, got %v", err)
	}
}

func TestValidateTagDTO_BlankNameFails(t *testing.T) {
	if err := validateTagDTO("   ", "#ffffff"); err == nil {
		t.Error("expected an error for a blank tag name")
	}
}

func TestValidateTagDTO_InvalidColorFails(t *testing.T) {
	if err := validateTagDTO("Urgent", "red"); err == nil {
		t.Error("expected an error for an invalid hex color")
	}
}

func TestIsValidHexColor_ValidLowercase(t *testing.T) {
	if !isValidHexColor("#abc123") {
		t.Error("expected #abc123 to be valid")
	}
}

func TestIsValidHexColor_ValidUppercase(t *testing.T) {
	if !isValidHexColor("#ABC123") {
		t.Error("expected #ABC123 to be valid")
	}
}

func TestIsValidHexColor_MissingHashIsInvalid(t *testing.T) {
	if isValidHexColor("abc123") {
		t.Error("expected a color without a leading # to be invalid")
	}
}

func TestIsValidHexColor_WrongLengthIsInvalid(t *testing.T) {
	if isValidHexColor("#abc") {
		t.Error("expected a 4-character color to be invalid")
	}
}

func TestIsValidHexColor_NonHexCharIsInvalid(t *testing.T) {
	if isValidHexColor("#gggggg") {
		t.Error("expected a color with non-hex characters to be invalid")
	}
}
