package invitation_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func GenerateInvitationCodeHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	var dto models.InvitationCodeCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	if dto.ExpiresInDays == 0 {
		dto.ExpiresInDays = 7
	}

	code, err := config.InvitationCodeService.CreateInvitationCode(
		r.Context(),
		userId,
		dto.ExpiresInDays,
	)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Invitation code generated successfully.",
		Data:    map[string]string{"code": code},
	})
}

func GetMyInvitationCodesHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	codes, err := config.InvitationCodeService.GetInvitationCodesByProfessional(
		r.Context(),
		userId,
	)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Invitation codes retrieved successfully.",
		Data:    codes,
	})
}

func ValidateInvitationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.InvitationCodeValidateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	if dto.Code == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Code is required."})
		return
	}

	invCode, err := config.InvitationCodeService.GetInvitationCodeByCode(
		r.Context(),
		dto.Code,
	)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Invitation code is valid.",
		Data:    map[string]int{"professional_id": invCode.ProfessionalID},
	})
}

func DeleteInvitationCodeHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	codeIDStr := chi.URLParam(r, "id")
	codeID, err := strconv.Atoi(codeIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid code ID format."})
		return
	}

	err = config.InvitationCodeService.DeleteInvitationCode(r.Context(), userId, codeID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Invitation code deleted successfully.",
		Data:    nil,
	})
}
